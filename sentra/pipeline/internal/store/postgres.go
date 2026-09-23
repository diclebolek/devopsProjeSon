package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

// Postgres stores telemetry in TimescaleDB when the extension is present,
// and in plain PostgreSQL otherwise. The table shape is the same either way.
type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, databaseURL string, logger *slog.Logger) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("postgres pool: %w", err)
	}
	store := &Postgres{pool: pool}
	if err := store.migrate(ctx, logger); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *Postgres) migrate(ctx context.Context, logger *slog.Logger) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS telemetry (
			captured_at TIMESTAMPTZ NOT NULL,
			device_id TEXT NOT NULL,
			schema_version INTEGER NOT NULL,
			firmware_version TEXT NOT NULL,
			sample_rate_hz INTEGER NOT NULL,
			window_samples INTEGER NOT NULL,
			rpm DOUBLE PRECISION NOT NULL,
			current_a DOUBLE PRECISION NOT NULL,
			voltage_v DOUBLE PRECISION NOT NULL,
			temperature_c DOUBLE PRECISION NOT NULL,
			accel_rms_g DOUBLE PRECISION NOT NULL,
			accel_peak_g DOUBLE PRECISION NOT NULL,
			crest_factor DOUBLE PRECISION NOT NULL,
			kurtosis DOUBLE PRECISION NOT NULL,
			dominant_freq_hz DOUBLE PRECISION NOT NULL,
			sensor_ok BOOLEAN NOT NULL,
			wifi_rssi_dbm INTEGER,
			ingested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (device_id, captured_at)
		)`,
		`CREATE INDEX IF NOT EXISTS telemetry_captured_at_idx ON telemetry (captured_at DESC)`,
		`CREATE TABLE IF NOT EXISTS condition_labels (
			device_id TEXT NOT NULL,
			captured_at TIMESTAMPTZ NOT NULL,
			condition TEXT NOT NULL,
			PRIMARY KEY (device_id, captured_at)
		)`,
		`CREATE TABLE IF NOT EXISTS device_status (
			device_id TEXT PRIMARY KEY,
			online BOOLEAN NOT NULL,
			firmware_version TEXT NOT NULL DEFAULT '',
			last_seen TIMESTAMPTZ NOT NULL,
			detail TEXT NOT NULL DEFAULT ''
		)`,
	}
	for _, statement := range statements {
		if _, err := p.pool.Exec(ctx, statement); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	if _, err := p.pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS timescaledb`); err != nil {
		logger.Warn("timescaledb extension is not available; telemetry stays a plain table", "error", err.Error())
		return nil
	}
	if _, err := p.pool.Exec(ctx, `SELECT create_hypertable('telemetry', 'captured_at', if_not_exists => TRUE)`); err != nil {
		logger.Warn("hypertable was not created; telemetry stays a plain table", "error", err.Error())
	}
	return nil
}

func (p *Postgres) InsertTelemetry(ctx context.Context, frame telemetry.Frame, ingestedAt time.Time) error {
	capturedAt := time.UnixMilli(frame.CapturedAtMs).UTC()
	tag, err := p.pool.Exec(ctx, `
		INSERT INTO telemetry (
			captured_at, device_id, schema_version, firmware_version, sample_rate_hz, window_samples,
			rpm, current_a, voltage_v, temperature_c, accel_rms_g, accel_peak_g, crest_factor,
			kurtosis, dominant_freq_hz, sensor_ok, wifi_rssi_dbm, ingested_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18
		)
		ON CONFLICT (device_id, captured_at) DO NOTHING`,
		capturedAt, frame.DeviceID, frame.SchemaVersion, frame.FirmwareVersion, frame.SampleRateHz, frame.WindowSamples,
		frame.RPM, frame.CurrentA, frame.VoltageV, frame.TemperatureC, frame.AccelRMS, frame.AccelPeak, frame.CrestFactor,
		frame.Kurtosis, frame.DominantFreqHz, frame.SensorOK, frame.WifiRssiDbm, ingestedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("insert telemetry: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return nil
	}
	_, err = p.pool.Exec(ctx, `
		INSERT INTO device_status (device_id, online, firmware_version, last_seen, detail)
		VALUES ($1, TRUE, $2, $3, '')
		ON CONFLICT (device_id) DO UPDATE SET
			online = TRUE,
			firmware_version = EXCLUDED.firmware_version,
			last_seen = EXCLUDED.last_seen`,
		frame.DeviceID, frame.FirmwareVersion, ingestedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("touch device status: %w", err)
	}
	return nil
}

func (p *Postgres) UpsertLabel(ctx context.Context, label telemetry.Label) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO condition_labels (device_id, captured_at, condition)
		VALUES ($1, $2, $3)
		ON CONFLICT (device_id, captured_at) DO UPDATE SET condition = EXCLUDED.condition`,
		label.DeviceID, time.UnixMilli(label.CapturedAtMs).UTC(), label.Condition,
	)
	if err != nil {
		return fmt.Errorf("upsert label: %w", err)
	}
	return nil
}

func (p *Postgres) UpsertStatus(ctx context.Context, deviceID string, status telemetry.Status, seenAt time.Time) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO device_status (device_id, online, firmware_version, last_seen, detail)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (device_id) DO UPDATE SET
			online = EXCLUDED.online,
			firmware_version = CASE WHEN EXCLUDED.firmware_version = '' THEN device_status.firmware_version ELSE EXCLUDED.firmware_version END,
			last_seen = EXCLUDED.last_seen,
			detail = EXCLUDED.detail`,
		deviceID, status.Online, status.FirmwareVersion, seenAt.UTC(), status.Detail,
	)
	if err != nil {
		return fmt.Errorf("upsert status: %w", err)
	}
	return nil
}

func (p *Postgres) ListTelemetry(ctx context.Context, deviceID string, limit int, _ time.Time) ([]telemetry.Reading, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT
			t.captured_at, t.device_id, t.schema_version, t.firmware_version, t.sample_rate_hz, t.window_samples,
			t.rpm, t.current_a, t.voltage_v, t.temperature_c, t.accel_rms_g, t.accel_peak_g, t.crest_factor,
			t.kurtosis, t.dominant_freq_hz, t.sensor_ok, t.wifi_rssi_dbm, t.ingested_at, c.condition
		FROM telemetry t
		LEFT JOIN condition_labels c
			ON c.device_id = t.device_id AND c.captured_at = t.captured_at
		WHERE ($1 = '' OR t.device_id = $1)
		ORDER BY t.captured_at DESC
		LIMIT $2`, deviceID, limit)
	if err != nil {
		return nil, fmt.Errorf("list telemetry: %w", err)
	}
	defer rows.Close()
	readings := make([]telemetry.Reading, 0, limit)
	for rows.Next() {
		reading, err := scanReading(rows)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list telemetry: %w", err)
	}
	return readings, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanReading(row scanner) (telemetry.Reading, error) {
	var (
		reading    telemetry.Reading
		capturedAt time.Time
		rssi       *int32
		condition  *string
	)
	err := row.Scan(
		&capturedAt, &reading.DeviceID, &reading.SchemaVersion, &reading.FirmwareVersion, &reading.SampleRateHz, &reading.WindowSamples,
		&reading.RPM, &reading.CurrentA, &reading.VoltageV, &reading.TemperatureC, &reading.AccelRMS, &reading.AccelPeak, &reading.CrestFactor,
		&reading.Kurtosis, &reading.DominantFreqHz, &reading.SensorOK, &rssi, &reading.IngestedAt, &condition,
	)
	if err != nil {
		return telemetry.Reading{}, fmt.Errorf("scan telemetry: %w", err)
	}
	reading.CapturedAtMs = capturedAt.UTC().UnixMilli()
	reading.IngestedAt = reading.IngestedAt.UTC()
	if rssi != nil {
		value := int(*rssi)
		reading.WifiRssiDbm = &value
	}
	if condition != nil {
		reading.ConditionLabel = *condition
	}
	return reading, nil
}

func (p *Postgres) ListDevices(ctx context.Context, now time.Time) ([]telemetry.Device, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT device_id, online, firmware_version, last_seen
		FROM device_status
		ORDER BY device_id`)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	defer rows.Close()
	var devices []telemetry.Device
	for rows.Next() {
		var device telemetry.Device
		if err := rows.Scan(&device.DeviceID, &device.Online, &device.FirmwareVersion, &device.LastSeen); err != nil {
			return nil, fmt.Errorf("scan device: %w", err)
		}
		device.LastSeen = device.LastSeen.UTC()
		device.Fresh = device.Online && now.Sub(device.LastSeen) <= FreshAfter
		devices = append(devices, device)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	for i := range devices {
		reading, err := p.latest(ctx, devices[i].DeviceID)
		if err != nil {
			return nil, err
		}
		if reading != nil {
			devices[i].Latest = reading
		}
	}
	return devices, nil
}

func (p *Postgres) latest(ctx context.Context, deviceID string) (*telemetry.Reading, error) {
	row := p.pool.QueryRow(ctx, `
		SELECT
			t.captured_at, t.device_id, t.schema_version, t.firmware_version, t.sample_rate_hz, t.window_samples,
			t.rpm, t.current_a, t.voltage_v, t.temperature_c, t.accel_rms_g, t.accel_peak_g, t.crest_factor,
			t.kurtosis, t.dominant_freq_hz, t.sensor_ok, t.wifi_rssi_dbm, t.ingested_at, c.condition
		FROM telemetry t
		LEFT JOIN condition_labels c
			ON c.device_id = t.device_id AND c.captured_at = t.captured_at
		WHERE t.device_id = $1
		ORDER BY t.captured_at DESC
		LIMIT 1`, deviceID)
	reading, err := scanReading(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reading, nil
}
