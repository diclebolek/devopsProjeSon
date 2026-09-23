package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/app"
	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/config"
	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/mqttpipe"
	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/store"
)

func main() {
	cfg, err := config.Load()
	logger := newLogger(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		logger.Error("configuration", "error", err.Error())
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	persistence, closeStore, err := openStore(ctx, cfg, logger)
	if err != nil {
		logger.Error("store", "error", err.Error())
		os.Exit(1)
	}
	if closeStore != nil {
		defer closeStore()
	}

	service := app.New(persistence, cfg.StoreMode, logger)
	if cfg.MQTTBroker != "" {
		stopMQTT := mqttpipe.Start(ctx, cfg.MQTTBroker, cfg.MQTTClientID, func(topic string, payload []byte) {
			service.HandleMQTT(context.Background(), topic, payload)
		}, service.SetMQTTConnected, logger)
		defer stopMQTT()
	} else {
		logger.Info("mqtt broker not set; accepting HTTP telemetry only")
	}

	logger.Info("listening", "addr", cfg.HTTPAddr, "store", cfg.StoreMode)
	if err := service.ListenAndServe(ctx, cfg.HTTPAddr); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("http server stopped", "error", err.Error())
		os.Exit(1)
	}
}

func openStore(ctx context.Context, cfg config.Config, logger *slog.Logger) (store.Store, func(), error) {
	if cfg.StoreMode == "memory" {
		logger.Info("using in-memory store")
		return store.NewMemory(), nil, nil
	}
	logger.Info("connecting postgres", "target", config.RedactedDatabaseTarget(cfg.DatabaseURL))
	var last error
	for attempt := 1; attempt <= 30; attempt++ {
		postgres, err := store.NewPostgres(ctx, cfg.DatabaseURL, logger)
		if err == nil {
			return postgres, postgres.Close, nil
		}
		last = err
		logger.Warn("postgres is not ready", "attempt", attempt, "error", err.Error())
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}
	return nil, nil, last
}

func newLogger(level, format string) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	options := &slog.HandlerOptions{Level: slogLevel}
	if format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stderr, options))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, options))
}
