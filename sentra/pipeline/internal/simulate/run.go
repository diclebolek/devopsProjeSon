package simulate

import (
	"context"
	"log/slog"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

// Publisher is how a frame leaves the simulator. MQTT is the real transport;
// HTTP is the same contract aimed at the ingest API when no broker is running.
type Publisher interface {
	PublishTelemetry(ctx context.Context, frame telemetry.Frame) error
	PublishLabel(ctx context.Context, label telemetry.Label) error
	PublishStatus(ctx context.Context, deviceID string, status telemetry.Status) error
}

// Config controls one simulated motor.
type Config struct {
	DeviceID        string
	FirmwareVersion string
	Regime          Regime
	RegimePeriod    int
	Interval        time.Duration
	MaxFrames       int
	Seed            int64
	Transport       string
	MQTTBroker      string
	IngestURL       string
}

// Run publishes status, then frames, until the context ends or MaxFrames is reached.
func Run(ctx context.Context, cfg Config, publisher Publisher, logger *slog.Logger) error {
	if cfg.RegimePeriod < 1 {
		cfg.RegimePeriod = 1
	}
	if err := publisher.PublishStatus(ctx, cfg.DeviceID, telemetry.Status{
		Online:          true,
		FirmwareVersion: cfg.FirmwareVersion,
		Detail:          "simulator",
	}); err != nil {
		logger.Warn("simulator could not announce online", "error", err.Error())
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := publisher.PublishStatus(stopCtx, cfg.DeviceID, telemetry.Status{
			Online:          false,
			FirmwareVersion: cfg.FirmwareVersion,
			Detail:          "simulator stopped",
		}); err != nil {
			logger.Warn("simulator could not announce offline", "error", err.Error())
		}
	}()

	for sequence := 0; cfg.MaxFrames <= 0 || sequence < cfg.MaxFrames; sequence++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		regime := Resolve(cfg.Regime, sequence, cfg.RegimePeriod)
		frame, label, err := Synthesize(regime, cfg.Seed, sequence, time.Now().UTC(), cfg.DeviceID, cfg.FirmwareVersion)
		if err != nil {
			return err
		}
		if err := publisher.PublishLabel(ctx, label); err != nil {
			logger.Warn("label publish failed", "error", err.Error(), "regime", string(regime))
		}
		if err := publisher.PublishTelemetry(ctx, frame); err != nil {
			logger.Warn("telemetry publish failed", "error", err.Error(), "regime", string(regime))
		} else {
			logger.Info("published frame", "device_id", frame.DeviceID, "regime", string(regime), "rpm", frame.RPM, "rms_g", frame.AccelRMS)
		}
		if cfg.MaxFrames > 0 && sequence+1 >= cfg.MaxFrames {
			return nil
		}
		timer := time.NewTimer(cfg.Interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return nil
}
