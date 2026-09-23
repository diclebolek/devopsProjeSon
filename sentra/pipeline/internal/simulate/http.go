package simulate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

// HTTPPublisher posts the same JSON the device would publish over MQTT.
type HTTPPublisher struct {
	base   string
	client *http.Client
}

func NewHTTPPublisher(baseURL string) *HTTPPublisher {
	return &HTTPPublisher{
		base: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p *HTTPPublisher) PublishTelemetry(ctx context.Context, frame telemetry.Frame) error {
	return p.post(ctx, "/api/v1/telemetry", frame)
}

func (p *HTTPPublisher) PublishLabel(ctx context.Context, label telemetry.Label) error {
	return p.post(ctx, "/api/v1/labels", label)
}

func (p *HTTPPublisher) PublishStatus(ctx context.Context, deviceID string, status telemetry.Status) error {
	body := telemetry.StatusReport{DeviceID: deviceID, Status: status}
	return p.post(ctx, "/api/v1/status", body)
}

func (p *HTTPPublisher) post(ctx context.Context, path string, body any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.base+path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("request %s: %w", path, err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := p.client.Do(request)
	if err != nil {
		return fmt.Errorf("post %s: %w", path, err)
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, response.Body)
	if response.StatusCode >= 300 {
		return fmt.Errorf("ingest %s returned %d", path, response.StatusCode)
	}
	return nil
}
