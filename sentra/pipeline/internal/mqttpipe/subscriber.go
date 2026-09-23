package mqttpipe

import (
	"context"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Handler receives one MQTT payload. It must not panic; the subscriber logs errors.
type Handler func(topic string, payload []byte)

// Start subscribes to the V1 topics and resubscribes after every reconnect.
// The returned function disconnects the client.
func Start(ctx context.Context, broker, clientID string, handler Handler, connected func(bool), logger *slog.Logger) func() {
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID(clientID)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectRetryInterval(2 * time.Second)
	options.SetKeepAlive(5 * time.Second)
	options.SetOrderMatters(false)
	options.SetCleanSession(true)
	options.SetOnConnectHandler(func(client mqtt.Client) {
		connected(true)
		logger.Info("mqtt connected", "broker", broker)
		for _, topic := range []string{"sentra/v1/+/telemetry", "sentra/v1/+/label", "sentra/v1/+/status"} {
			token := client.Subscribe(topic, 1, func(_ mqtt.Client, message mqtt.Message) {
				handler(message.Topic(), message.Payload())
			})
			if token.WaitTimeout(5*time.Second) && token.Error() != nil {
				logger.Error("mqtt subscribe failed", "topic", topic, "error", token.Error().Error())
			}
		}
	})
	options.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		connected(false)
		logger.Warn("mqtt connection lost", "error", err.Error())
	})

	client := mqtt.NewClient(options)
	go func() {
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			logger.Error("mqtt connect failed", "error", token.Error().Error())
		}
	}()

	return func() {
		done := make(chan struct{})
		go func() {
			client.Disconnect(500)
			close(done)
		}()
		select {
		case <-done:
		case <-ctx.Done():
		case <-time.After(2 * time.Second):
		}
		connected(false)
	}
}
