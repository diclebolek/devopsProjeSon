package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/simulate"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg, err := simulate.LoadEnv()
	if err != nil {
		logger.Error("configuration", "error", err.Error())
		os.Exit(2)
	}

	var publisher simulate.Publisher
	switch cfg.Transport {
	case "http":
		publisher = simulate.NewHTTPPublisher(cfg.IngestURL)
		logger.Info("simulator transport http", "url", cfg.IngestURL)
	case "mqtt":
		publisher, err = simulate.DialMQTT(cfg.MQTTBroker, "sentra-simulator-"+cfg.DeviceID, cfg.DeviceID, cfg.FirmwareVersion, logger)
		if err != nil {
			logger.Error("mqtt", "error", err.Error())
			os.Exit(1)
		}
		logger.Info("simulator transport mqtt", "broker", cfg.MQTTBroker)
	default:
		logger.Error("unknown transport", "transport", cfg.Transport)
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := simulate.Run(ctx, cfg, publisher, logger); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("simulator stopped", "error", err.Error())
		os.Exit(1)
	}
}
