package telemetry

import (
	"math"
	"regexp"
)

const (
	minSampleRateHz = 100
	maxSampleRateHz = 8000
	minWindow       = 16
	maxWindow       = 2048
)

var (
	deviceIDPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$`)
	versionPattern  = regexp.MustCompile(`^[0-9A-Za-z._-]{1,32}$`)
	conditions      = map[string]struct{}{
		"NORMAL":         {},
		"OVERLOAD":       {},
		"UNBALANCED":     {},
		"HIGH_VIBRATION": {},
	}
)

// AllowedConditions is the closed set used by labels. Models in later stages
// must emit one of these names or stay silent.
// ValidDeviceID reports whether a device slug matches the wire contract.
func ValidDeviceID(deviceID string) bool {
	return deviceIDPattern.MatchString(deviceID)
}

func AllowedConditions() []string {
	return []string{"NORMAL", "OVERLOAD", "UNBALANCED", "HIGH_VIBRATION"}
}

func (f Frame) Validate() error {
	if f.SchemaVersion != SchemaVersion {
		return invalid("schema_version", "unsupported schema_version")
	}
	if !deviceIDPattern.MatchString(f.DeviceID) {
		return invalid("device_id", "use a lowercase slug, 1-64 characters")
	}
	if !versionPattern.MatchString(f.FirmwareVersion) {
		return invalid("firmware_version", "empty or longer than 32 characters")
	}
	if f.CapturedAtMs <= 0 {
		return invalid("captured_at_ms", "must be a positive unix millisecond timestamp")
	}
	if f.SampleRateHz < minSampleRateHz || f.SampleRateHz > maxSampleRateHz {
		return invalid("sample_rate_hz", "must be between 100 and 8000")
	}
	if f.WindowSamples < minWindow || f.WindowSamples > maxWindow || !isPowerOfTwo(f.WindowSamples) {
		return invalid("window_samples", "must be a power of two between 16 and 2048")
	}
	if err := finiteRange(f.RPM, 0, 20_000, "rpm"); err != nil {
		return err
	}
	if err := finiteRange(f.CurrentA, 0, 20, "current_a"); err != nil {
		return err
	}
	if err := finiteRange(f.VoltageV, 0, 30, "voltage_v"); err != nil {
		return err
	}
	if err := finiteRange(f.TemperatureC, -40, 150, "temperature_c"); err != nil {
		return err
	}
	if err := finiteRange(f.AccelRMS, 0, 50, "accel_rms_g"); err != nil {
		return err
	}
	if err := finiteRange(f.AccelPeak, 0, 50, "accel_peak_g"); err != nil {
		return err
	}
	if err := finiteRange(f.CrestFactor, 0, 100, "crest_factor"); err != nil {
		return err
	}
	if err := finiteRange(f.Kurtosis, 0, 100, "kurtosis"); err != nil {
		return err
	}
	nyquist := float64(f.SampleRateHz) / 2
	if err := finiteRange(f.DominantFreqHz, 0, nyquist, "dominant_freq_hz"); err != nil {
		return err
	}
	if f.WifiRssiDbm != nil && (*f.WifiRssiDbm < -120 || *f.WifiRssiDbm > 0) {
		return invalid("wifi_rssi_dbm", "must be between -120 and 0")
	}
	return nil
}

func (l Label) Validate() error {
	if l.SchemaVersion != SchemaVersion {
		return invalid("schema_version", "unsupported schema_version")
	}
	if !deviceIDPattern.MatchString(l.DeviceID) {
		return invalid("device_id", "use a lowercase slug, 1-64 characters")
	}
	if l.CapturedAtMs <= 0 {
		return invalid("captured_at_ms", "must be a positive unix millisecond timestamp")
	}
	if _, ok := conditions[l.Condition]; !ok {
		return invalid("condition", "must be NORMAL, OVERLOAD, UNBALANCED, or HIGH_VIBRATION")
	}
	return nil
}

func ValidateStatus(deviceID string, status Status) error {
	if !deviceIDPattern.MatchString(deviceID) {
		return invalid("device_id", "use a lowercase slug, 1-64 characters")
	}
	if status.FirmwareVersion != "" && !versionPattern.MatchString(status.FirmwareVersion) {
		return invalid("firmware_version", "empty or longer than 32 characters")
	}
	if len(status.Detail) > 200 {
		return invalid("detail", "longer than 200 characters")
	}
	return nil
}

func finiteRange(value, min, max float64, field string) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return invalid(field, "must be a finite number")
	}
	if value < min || value > max {
		return invalid(field, "outside the accepted range")
	}
	return nil
}

func invalid(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}

func isPowerOfTwo(value int) bool {
	return value > 0 && value&(value-1) == 0
}
