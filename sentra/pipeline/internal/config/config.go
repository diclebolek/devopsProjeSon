package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Config is the ingest process configuration. Secrets stay in the environment.
type Config struct {
	HTTPAddr     string
	StoreMode    string
	DatabaseURL  string
	MQTTBroker   string
	MQTTClientID string
	LogLevel     string
	LogFormat    string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:     envOr("SENTRA_HTTP_ADDR", ":8080"),
		StoreMode:    strings.ToLower(envOr("SENTRA_STORE", "memory")),
		DatabaseURL:  os.Getenv("SENTRA_DATABASE_URL"),
		MQTTBroker:   os.Getenv("SENTRA_MQTT_BROKER"),
		MQTTClientID: envOr("SENTRA_MQTT_CLIENT_ID", "sentra-ingest"),
		LogLevel:     strings.ToLower(envOr("SENTRA_LOG_LEVEL", "info")),
		LogFormat:    strings.ToLower(envOr("SENTRA_LOG_FORMAT", "text")),
	}
	switch cfg.StoreMode {
	case "memory", "postgres":
	default:
		return Config{}, fmt.Errorf("SENTRA_STORE must be memory or postgres")
	}
	if cfg.StoreMode == "postgres" && strings.TrimSpace(cfg.DatabaseURL) == "" {
		return Config{}, fmt.Errorf("SENTRA_DATABASE_URL is required when SENTRA_STORE=postgres")
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

// RedactedDatabaseTarget returns host and database name without credentials.
func RedactedDatabaseTarget(databaseURL string) string {
	parsed, err := url.Parse(databaseURL)
	if err != nil || parsed.Host == "" {
		return "unparsed"
	}
	return parsed.Host + parsed.Path
}
