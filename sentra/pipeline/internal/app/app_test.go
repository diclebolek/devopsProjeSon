package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/store"
	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

func TestDashboardAndIngestRoundTrip(t *testing.T) {
	svc := New(store.NewMemory(), "memory", discardLogger())
	fixed := time.UnixMilli(1_710_000_000_000).UTC()
	svc.now = func() time.Time { return fixed }
	server := httptest.NewServer(svc.Handler())
	defer server.Close()

	page, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer page.Body.Close()
	body, _ := readAll(page)
	if page.StatusCode != http.StatusOK || !strings.Contains(body, "SENTRA") || !strings.Contains(body, "Kenar hesaplamalı") {
		t.Fatalf("dashboard status %d body %s", page.StatusCode, body[:smaller(120, len(body))])
	}

	label := telemetry.Label{
		SchemaVersion: 1,
		DeviceID:      "motor-lab-01",
		CapturedAtMs:  1_710_000_000_000,
		Condition:     "OVERLOAD",
	}
	postJSON(t, server.URL+"/api/v1/labels", label, http.StatusCreated)

	raw, err := os.ReadFile("../../../testdata/frame_normal.json")
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.Post(server.URL+"/api/v1/telemetry", "application/json", bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("telemetry status %d", response.StatusCode)
	}

	listResponse, err := http.Get(server.URL + "/api/v1/telemetry?device_id=motor-lab-01&limit=10")
	if err != nil {
		t.Fatal(err)
	}
	defer listResponse.Body.Close()
	var listed struct {
		Items []telemetry.Reading `json:"items"`
		Count int                 `json:"count"`
	}
	if err := json.NewDecoder(listResponse.Body).Decode(&listed); err != nil {
		t.Fatal(err)
	}
	if listed.Count != 1 || listed.Items[0].ConditionLabel != "OVERLOAD" {
		t.Fatalf("list %#v", listed)
	}
	if listed.Items[0].RPM != 1500 {
		t.Fatalf("rpm %v", listed.Items[0].RPM)
	}

	devicesResponse, err := http.Get(server.URL + "/api/v1/devices")
	if err != nil {
		t.Fatal(err)
	}
	defer devicesResponse.Body.Close()
	var devices struct {
		Devices []telemetry.Device `json:"devices"`
	}
	if err := json.NewDecoder(devicesResponse.Body).Decode(&devices); err != nil {
		t.Fatal(err)
	}
	if len(devices.Devices) != 1 || !devices.Devices[0].Fresh || devices.Devices[0].Latest == nil {
		t.Fatalf("devices %#v", devices.Devices)
	}

	svc.now = func() time.Time { return fixed.Add(20 * time.Second) }
	devicesResponse, err = http.Get(server.URL + "/api/v1/devices")
	if err != nil {
		t.Fatal(err)
	}
	defer devicesResponse.Body.Close()
	if err := json.NewDecoder(devicesResponse.Body).Decode(&devices); err != nil {
		t.Fatal(err)
	}
	if devices.Devices[0].Fresh {
		t.Fatal("stale device still marked fresh")
	}
}

func TestRejectsBadPayloadsAndTopicMismatch(t *testing.T) {
	svc := New(store.NewMemory(), "memory", discardLogger())
	server := httptest.NewServer(svc.Handler())
	defer server.Close()

	response, err := http.Post(server.URL+"/api/v1/telemetry", "application/json", strings.NewReader(`{"device_id":"NOPE"}`))
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("status %d", response.StatusCode)
	}

	plain, err := http.Post(server.URL+"/api/v1/telemetry", "text/plain", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	plain.Body.Close()
	if plain.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("content type status %d", plain.StatusCode)
	}

	badDevice, err := http.Get(server.URL + "/api/v1/telemetry?device_id=Bad_ID")
	if err != nil {
		t.Fatal(err)
	}
	badDevice.Body.Close()
	if badDevice.StatusCode != http.StatusBadRequest {
		t.Fatalf("query status %d", badDevice.StatusCode)
	}

	svc.HandleMQTT(context.Background(), "sentra/v1/motor-lab-02/telemetry", []byte(`{"schema_version":1,"device_id":"motor-lab-01"}`))
	if svc.Metrics.Rejected.Load() == 0 {
		t.Fatal("expected mqtt rejection")
	}

	health, err := http.Get(server.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer health.Body.Close()
	healthBody, _ := readAll(health)
	if !strings.Contains(healthBody, `"store":"memory"`) {
		t.Fatal(healthBody)
	}
	metrics, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer metrics.Body.Close()
	metricsBody, _ := readAll(metrics)
	if !strings.Contains(metricsBody, "sentra_messages_rejected_total") {
		t.Fatal(metricsBody)
	}
}

func postJSON(t *testing.T, url string, body any, want int) {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != want {
		t.Fatalf("POST %s status %d", url, response.StatusCode)
	}
}

func readAll(response *http.Response) (string, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(response.Body)
	return buf.String(), err
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func smaller(a, b int) int {
	if a < b {
		return a
	}
	return b
}
