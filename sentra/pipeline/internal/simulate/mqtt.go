package simulate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

// MQTTPublisher emits the V1 topics. Connection retries are handled by the client.
type MQTTPublisher struct {
	client mqtt.Client
}

func DialMQTT(broker, clientID, deviceID, firmware string, logger *slog.Logger) (*MQTTPublisher, error) {
	if broker == "" {
		return nil, errors.New("mqtt broker is empty")
	}
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID(clientID)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectRetryInterval(2 * time.Second)
	options.SetKeepAlive(5 * time.Second)
	will := fmt.Sprintf(`{"online":false,"firmware_version":%q,"detail":"connection lost"}`, firmware)
	options.SetWill(fmt.Sprintf("sentra/v1/%s/status", deviceID), will, 1, true)
	options.SetOnConnectHandler(func(mqtt.Client) {
		logger.Info("simulator mqtt connected", "broker", broker)
	})
	options.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logger.Warn("simulator mqtt connection lost", "error", err.Error())
	})
	client := mqtt.NewClient(options)
	token := client.Connect()
	go func() {
		if token.Wait() && token.Error() != nil {
			logger.Error("simulator mqtt connect failed", "error", token.Error().Error())
		}
	}()
	return &MQTTPublisher{client: client}, nil
}

func (p *MQTTPublisher) PublishTelemetry(ctx context.Context, frame telemetry.Frame) error {
	return p.publish(ctx, fmt.Sprintf("sentra/v1/%s/telemetry", frame.DeviceID), false, frame)
}

func (p *MQTTPublisher) PublishLabel(ctx context.Context, label telemetry.Label) error {
	return p.publish(ctx, fmt.Sprintf("sentra/v1/%s/label", label.DeviceID), false, label)
}

func (p *MQTTPublisher) PublishStatus(ctx context.Context, deviceID string, status telemetry.Status) error {
	return p.publish(ctx, fmt.Sprintf("sentra/v1/%s/status", deviceID), true, status)
}

func (p *MQTTPublisher) publish(ctx context.Context, topic string, retained bool, body any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !p.client.IsConnected() {
		return errors.New("mqtt is not connected")
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encode mqtt payload: %w", err)
	}
	token := p.client.Publish(topic, 1, retained, payload)
	if !token.WaitTimeout(3 * time.Second) {
		return errors.New("mqtt publish timed out")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt publish: %w", err)
	}
	return nil
}
