package simulate

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

// LoadEnv reads the simulator process configuration.
func LoadEnv() (Config, error) {
	regime, err := ParseRegime(strings.ToUpper(strings.TrimSpace(os.Getenv("SENTRA_REGIME"))))
	if err != nil {
		return Config{}, err
	}
	period, err := envInt("SENTRA_REGIME_PERIOD", 8)
	if err != nil {
		return Config{}, err
	}
	intervalMS, err := envInt("SENTRA_INTERVAL_MS", 1000)
	if err != nil {
		return Config{}, err
	}
	maxFrames, err := envInt("SENTRA_MAX_FRAMES", 0)
	if err != nil {
		return Config{}, err
	}
	seed, err := envInt("SENTRA_SEED", 1)
	if err != nil {
		return Config{}, err
	}
	transport := strings.ToLower(envOr("SENTRA_TRANSPORT", "mqtt"))
	if transport != "mqtt" && transport != "http" {
		return Config{}, fmt.Errorf("SENTRA_TRANSPORT must be mqtt or http")
	}
	cfg := Config{
		DeviceID:        envOr("SENTRA_DEVICE_ID", "motor-lab-01"),
		FirmwareVersion: envOr("SENTRA_FIRMWARE_VERSION", "0.1.0"),
		Regime:          regime,
		RegimePeriod:    period,
		Interval:        time.Duration(intervalMS) * time.Millisecond,
		MaxFrames:       maxFrames,
		Seed:            int64(seed),
		Transport:       transport,
		MQTTBroker:      envOr("SENTRA_MQTT_BROKER", "tcp://127.0.0.1:1883"),
		IngestURL:       envOr("SENTRA_INGEST_URL", "http://127.0.0.1:8080"),
	}
	if !telemetry.ValidDeviceID(cfg.DeviceID) {
		return Config{}, fmt.Errorf("SENTRA_DEVICE_ID is not a lowercase slug")
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", key)
	}
	return value, nil
}
