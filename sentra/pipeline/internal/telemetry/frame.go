package telemetry

import "time"

// SchemaVersion is the only payload generation this binary accepts.
const SchemaVersion = 1

// Frame is the compact edge payload. Field names are the wire contract.
type Frame struct {
	SchemaVersion   int     `json:"schema_version"`
	DeviceID        string  `json:"device_id"`
	FirmwareVersion string  `json:"firmware_version"`
	CapturedAtMs    int64   `json:"captured_at_ms"`
	SampleRateHz    int     `json:"sample_rate_hz"`
	WindowSamples   int     `json:"window_samples"`
	RPM             float64 `json:"rpm"`
	CurrentA        float64 `json:"current_a"`
	VoltageV        float64 `json:"voltage_v"`
	TemperatureC    float64 `json:"temperature_c"`
	AccelRMS        float64 `json:"accel_rms_g"`
	AccelPeak       float64 `json:"accel_peak_g"`
	CrestFactor     float64 `json:"crest_factor"`
	Kurtosis        float64 `json:"kurtosis"`
	DominantFreqHz  float64 `json:"dominant_freq_hz"`
	SensorOK        bool    `json:"sensor_ok"`
	WifiRssiDbm     *int    `json:"wifi_rssi_dbm,omitempty"`
}

// Reading is a stored frame plus the simulator or technician label, if one arrived.
type Reading struct {
	Frame
	ConditionLabel string    `json:"condition_label,omitempty"`
	IngestedAt     time.Time `json:"ingested_at"`
}

// Label is ground truth from the simulator or a later labeling tool.
// It is not a model prediction.
type Label struct {
	SchemaVersion int    `json:"schema_version"`
	DeviceID      string `json:"device_id"`
	CapturedAtMs  int64  `json:"captured_at_ms"`
	Condition     string `json:"condition"`
}

// Status is the retained online/offline announcement, including MQTT last-will.
type Status struct {
	Online          bool   `json:"online"`
	FirmwareVersion string `json:"firmware_version"`
	Detail          string `json:"detail"`
}

// StatusReport carries a device id when there is no MQTT topic to take it from.
type StatusReport struct {
	DeviceID string `json:"device_id"`
	Status
}

// Device is the latest known state of one motor node.
type Device struct {
	DeviceID        string    `json:"device_id"`
	Online          bool      `json:"online"`
	FirmwareVersion string    `json:"firmware_version,omitempty"`
	LastSeen        time.Time `json:"last_seen"`
	Fresh           bool      `json:"fresh"`
	Latest          *Reading  `json:"latest,omitempty"`
}

// ValidationError is a client payload problem. Callers map it to HTTP 400
// or a rejected MQTT counter. It is not a server failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return e.Field + ": " + e.Message
}
