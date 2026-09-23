package simulate

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

func TestHTTPPublisherUsesTheIngestPaths(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		paths = append(paths, r.URL.Path+" "+string(body))
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("content type %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	publisher := NewHTTPPublisher(server.URL)
	capturedAt := time.UnixMilli(1_710_000_000_000).UTC()
	frame, label, err := Synthesize(RegimeNormal, 1, 0, capturedAt, "motor-lab-01", "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if err := publisher.PublishLabel(context.Background(), label); err != nil {
		t.Fatal(err)
	}
	if err := publisher.PublishTelemetry(context.Background(), frame); err != nil {
		t.Fatal(err)
	}
	if err := publisher.PublishStatus(context.Background(), "motor-lab-01", telemetry.Status{
		Online:          true,
		FirmwareVersion: "0.1.0",
		Detail:          "simulator",
	}); err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(paths, "\n")
	for _, path := range []string{"/api/v1/labels", "/api/v1/telemetry", "/api/v1/status"} {
		if !strings.Contains(joined, path) {
			t.Fatalf("missing %s in %s", path, joined)
		}
	}
}
