package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/store"
	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/web"
)

const maxPayloadBytes = 64 * 1024

// Metrics counts accepted and rejected messages for /metrics.
type Metrics struct {
	Received atomic.Uint64
	Stored   atomic.Uint64
	Rejected atomic.Uint64
}

// Service is the ingest application: HTTP, MQTT, and the dashboard share it.
type Service struct {
	store     store.Store
	storeName string
	log       *slog.Logger
	Metrics   Metrics
	now       func() time.Time
	mqttUp    atomic.Bool
}

func New(st store.Store, storeName string, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		store:     st,
		storeName: storeName,
		log:       logger,
		now:       time.Now,
	}
}

func (s *Service) SetMQTTConnected(connected bool) {
	s.mqttUp.Store(connected)
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /readyz", s.ready)
	mux.HandleFunc("GET /metrics", s.metrics)
	mux.HandleFunc("GET /{$}", s.dashboard)
	mux.HandleFunc("GET /api/v1/telemetry", s.listTelemetry)
	mux.HandleFunc("GET /api/v1/devices", s.listDevices)
	mux.HandleFunc("POST /api/v1/telemetry", s.postTelemetry)
	mux.HandleFunc("POST /api/v1/labels", s.postLabel)
	mux.HandleFunc("POST /api/v1/status", s.postStatus)
	return s.recover(mux)
}

func (s *Service) ListenAndServe(ctx context.Context, addr string) error {
	server := &http.Server{
		Addr:              addr,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown http: %w", err)
		}
		return nil
	case err := <-errCh:
		return err
	}
}

// HandleMQTT is the broker callback. Bad payloads are counted and logged; they
// do not stop the subscriber.
func (s *Service) HandleMQTT(ctx context.Context, topic string, payload []byte) {
	s.Metrics.Received.Add(1)
	if len(payload) > maxPayloadBytes {
		s.reject(topic, errors.New("payload is larger than 64 KiB"))
		return
	}
	deviceID, kind, err := telemetry.ParseTopic(topic)
	if err != nil {
		s.reject(topic, err)
		return
	}
	switch kind {
	case "telemetry":
		frame, err := telemetry.DecodeFrame(payload)
		if err != nil {
			s.reject(topic, err)
			return
		}
		if err := telemetry.SameDevice(deviceID, frame.DeviceID); err != nil {
			s.reject(topic, err)
			return
		}
		s.finish(topic, s.acceptFrame(ctx, frame))
	case "label":
		label, err := telemetry.DecodeLabel(payload)
		if err != nil {
			s.reject(topic, err)
			return
		}
		if err := telemetry.SameDevice(deviceID, label.DeviceID); err != nil {
			s.reject(topic, err)
			return
		}
		s.finish(topic, s.acceptLabel(ctx, label))
	case "status":
		status, err := telemetry.DecodeStatus(payload)
		if err != nil {
			s.reject(topic, err)
			return
		}
		s.finish(topic, s.acceptStatus(ctx, deviceID, status))
	default:
		s.reject(topic, errors.New("unsupported topic"))
	}
}

func (s *Service) finish(topic string, err error) {
	if err == nil {
		s.Metrics.Stored.Add(1)
		return
	}
	s.reject(topic, err)
}

func (s *Service) reject(topic string, err error) {
	s.Metrics.Rejected.Add(1)
	s.log.Warn("rejected message", "topic", topic, "error", err.Error())
}

func (s *Service) acceptFrame(ctx context.Context, frame telemetry.Frame) error {
	if err := frame.Validate(); err != nil {
		return err
	}
	return s.store.InsertTelemetry(ctx, frame, s.now().UTC())
}

func (s *Service) acceptLabel(ctx context.Context, label telemetry.Label) error {
	if err := label.Validate(); err != nil {
		return err
	}
	return s.store.UpsertLabel(ctx, label)
}

func (s *Service) acceptStatus(ctx context.Context, deviceID string, status telemetry.Status) error {
	if err := telemetry.ValidateStatus(deviceID, status); err != nil {
		return err
	}
	return s.store.UpsertStatus(ctx, deviceID, status, s.now().UTC())
}

func (s *Service) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":         "ok",
		"store":          s.storeName,
		"mqtt_connected": s.mqttUp.Load(),
	})
}

func (s *Service) ready(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Ping(r.Context()); err != nil {
		s.log.Error("readiness ping failed", "error", err.Error())
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Service) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = fmt.Fprintf(w, `# HELP sentra_messages_received_total Messages observed on MQTT.
# TYPE sentra_messages_received_total counter
sentra_messages_received_total %d
# HELP sentra_messages_stored_total Messages accepted into the store.
# TYPE sentra_messages_stored_total counter
sentra_messages_stored_total %d
# HELP sentra_messages_rejected_total Messages rejected by validation.
# TYPE sentra_messages_rejected_total counter
sentra_messages_rejected_total %d
`, s.Metrics.Received.Load(), s.Metrics.Stored.Load(), s.Metrics.Rejected.Load())
}

func (s *Service) dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(web.IndexHTML)
}

func (s *Service) listTelemetry(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device_id")
	if deviceID != "" && !telemetryValidID(deviceID) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "device_id is not a valid slug"})
		return
	}
	limit, err := parseLimit(r.URL.Query().Get("limit"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	items, err := s.store.ListTelemetry(r.Context(), deviceID, limit, s.now().UTC())
	if err != nil {
		s.log.Error("list telemetry", "error", err.Error())
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read telemetry"})
		return
	}
	if items == nil {
		items = []telemetry.Reading{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "count": len(items)})
}

func (s *Service) listDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.store.ListDevices(r.Context(), s.now().UTC())
	if err != nil {
		s.log.Error("list devices", "error", err.Error())
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read devices"})
		return
	}
	if devices == nil {
		devices = []telemetry.Device{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"devices": devices, "count": len(devices)})
}

func (s *Service) postTelemetry(w http.ResponseWriter, r *http.Request) {
	var frame telemetry.Frame
	if !readJSON(w, r, &frame) {
		return
	}
	if err := s.acceptFrame(r.Context(), frame); err != nil {
		s.writeFailure(w, err)
		return
	}
	s.Metrics.Stored.Add(1)
	writeJSON(w, http.StatusCreated, map[string]string{"status": "stored"})
}

func (s *Service) postLabel(w http.ResponseWriter, r *http.Request) {
	var label telemetry.Label
	if !readJSON(w, r, &label) {
		return
	}
	if err := s.acceptLabel(r.Context(), label); err != nil {
		s.writeFailure(w, err)
		return
	}
	s.Metrics.Stored.Add(1)
	writeJSON(w, http.StatusCreated, map[string]string{"status": "stored"})
}

func (s *Service) postStatus(w http.ResponseWriter, r *http.Request) {
	var report telemetry.StatusReport
	if !readJSON(w, r, &report) {
		return
	}
	if err := s.acceptStatus(r.Context(), report.DeviceID, report.Status); err != nil {
		s.writeFailure(w, err)
		return
	}
	s.Metrics.Stored.Add(1)
	writeJSON(w, http.StatusCreated, map[string]string{"status": "stored"})
}

func (s *Service) writeFailure(w http.ResponseWriter, err error) {
	var validation *telemetry.ValidationError
	if errors.As(err, &validation) {
		s.Metrics.Rejected.Add(1)
		s.log.Warn("rejected message", "error", err.Error())
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	s.log.Error("store failure", "error", err.Error())
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not store message"})
}

func (s *Service) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				s.log.Error("panic", "error", fmt.Sprint(recovered), "path", r.URL.Path)
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && contentType != "application/json" && !hasJSONPrefix(contentType) {
		writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{"error": "content type must be application/json"})
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxPayloadBytes)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "payload is larger than 64 KiB"})
			return false
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body is not valid json"})
		return false
	}
	return true
}

func hasJSONPrefix(contentType string) bool {
	return len(contentType) >= len("application/json") && contentType[:len("application/json")] == "application/json"
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return
	}
}

func parseLimit(raw string) (int, error) {
	if raw == "" {
		return 50, nil
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 || limit > 200 {
		return 0, errors.New("limit must be between 1 and 200")
	}
	return limit, nil
}

func telemetryValidID(deviceID string) bool {
	return telemetry.ValidDeviceID(deviceID)
}
