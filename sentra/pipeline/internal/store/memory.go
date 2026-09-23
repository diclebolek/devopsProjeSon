package store

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

const memoryCap = 2000

type memoryStatus struct {
	online          bool
	firmwareVersion string
	detail          string
	lastSeen        time.Time
}

// Memory is a process-local store for tests and for running the dashboard without Docker.
type Memory struct {
	mu       sync.RWMutex
	readings []telemetry.Reading
	labels   map[string]string
	statuses map[string]memoryStatus
	latest   map[string]telemetry.Reading
}

func NewMemory() *Memory {
	return &Memory{
		labels:   map[string]string{},
		statuses: map[string]memoryStatus{},
		latest:   map[string]telemetry.Reading{},
	}
}

func (m *Memory) InsertTelemetry(_ context.Context, frame telemetry.Frame, ingestedAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	reading := telemetry.Reading{
		Frame:          frame,
		ConditionLabel: m.labels[labelKey(frame.DeviceID, frame.CapturedAtMs)],
		IngestedAt:     ingestedAt.UTC(),
	}
	m.readings = append(m.readings, reading)
	if len(m.readings) > memoryCap {
		m.readings = m.readings[len(m.readings)-memoryCap:]
	}
	m.latest[frame.DeviceID] = reading
	current := m.statuses[frame.DeviceID]
	current.online = true
	current.firmwareVersion = frame.FirmwareVersion
	current.lastSeen = ingestedAt.UTC()
	m.statuses[frame.DeviceID] = current
	return nil
}

func (m *Memory) UpsertLabel(_ context.Context, label telemetry.Label) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.labels[labelKey(label.DeviceID, label.CapturedAtMs)] = label.Condition
	for i := range m.readings {
		if m.readings[i].DeviceID == label.DeviceID && m.readings[i].CapturedAtMs == label.CapturedAtMs {
			m.readings[i].ConditionLabel = label.Condition
		}
	}
	if latest, ok := m.latest[label.DeviceID]; ok && latest.CapturedAtMs == label.CapturedAtMs {
		latest.ConditionLabel = label.Condition
		m.latest[label.DeviceID] = latest
	}
	return nil
}

func (m *Memory) UpsertStatus(_ context.Context, deviceID string, status telemetry.Status, seenAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	current := m.statuses[deviceID]
	current.online = status.Online
	if status.FirmwareVersion != "" {
		current.firmwareVersion = status.FirmwareVersion
	}
	current.detail = status.Detail
	current.lastSeen = seenAt.UTC()
	m.statuses[deviceID] = current
	return nil
}

func (m *Memory) ListTelemetry(_ context.Context, deviceID string, limit int, _ time.Time) ([]telemetry.Reading, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]telemetry.Reading, 0, limit)
	for i := len(m.readings) - 1; i >= 0 && len(out) < limit; i-- {
		if deviceID == "" || m.readings[i].DeviceID == deviceID {
			reading := m.readings[i]
			if reading.ConditionLabel == "" {
				reading.ConditionLabel = m.labels[labelKey(reading.DeviceID, reading.CapturedAtMs)]
			}
			out = append(out, reading)
		}
	}
	return out, nil
}

func (m *Memory) ListDevices(_ context.Context, now time.Time) ([]telemetry.Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	devices := make([]telemetry.Device, 0, len(m.statuses))
	for id, status := range m.statuses {
		device := telemetry.Device{
			DeviceID:        id,
			Online:          status.online,
			FirmwareVersion: status.firmwareVersion,
			LastSeen:        status.lastSeen,
			Fresh:           status.online && now.Sub(status.lastSeen) <= FreshAfter,
		}
		if latest, ok := m.latest[id]; ok {
			copy := latest
			device.Latest = &copy
		}
		devices = append(devices, device)
	}
	sort.Slice(devices, func(i, j int) bool { return devices[i].DeviceID < devices[j].DeviceID })
	return devices, nil
}

func (m *Memory) Ping(context.Context) error { return nil }

func labelKey(deviceID string, capturedAtMs int64) string {
	return deviceID + "|" + strconv.FormatInt(capturedAtMs, 10)
}
