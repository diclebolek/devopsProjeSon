package simulate

import (
	"fmt"
	"math"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/dsp"
	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

const (
	sampleRateHz  = 1000
	windowSamples = 256
)

// Regime is a named operating condition. CYCLE rotates the concrete regimes.
type Regime string

const (
	RegimeNormal        Regime = "NORMAL"
	RegimeOverload      Regime = "OVERLOAD"
	RegimeUnbalanced    Regime = "UNBALANCED"
	RegimeHighVibration Regime = "HIGH_VIBRATION"
	RegimeCycle         Regime = "CYCLE"
)

func ParseRegime(raw string) (Regime, error) {
	switch Regime(raw) {
	case RegimeNormal, RegimeOverload, RegimeUnbalanced, RegimeHighVibration, RegimeCycle, "":
		if raw == "" {
			return RegimeCycle, nil
		}
		return Regime(raw), nil
	default:
		return "", fmt.Errorf("unknown regime %q", raw)
	}
}

// Resolve picks the concrete regime for this frame. Period is measured in frames.
func Resolve(requested Regime, sequence, period int) Regime {
	if requested != RegimeCycle {
		return requested
	}
	if period < 1 {
		period = 1
	}
	order := []Regime{RegimeNormal, RegimeOverload, RegimeUnbalanced, RegimeHighVibration}
	if sequence < 0 {
		sequence = 0
	}
	return order[(sequence/period)%len(order)]
}

type profile struct {
	rpm       float64
	currentA  float64
	voltageV  float64
	tempC     float64
	amplitude float64
	impulsive bool
}

func profileFor(regime Regime) profile {
	switch regime {
	case RegimeOverload:
		return profile{rpm: 1180, currentA: 2.45, voltageV: 11.15, tempC: 72, amplitude: 0.22}
	case RegimeUnbalanced:
		return profile{rpm: 1490, currentA: 0.98, voltageV: 11.9, tempC: 41, amplitude: 0.85}
	case RegimeHighVibration:
		return profile{rpm: 1510, currentA: 0.9, voltageV: 12.0, tempC: 44, amplitude: 0.25, impulsive: true}
	default:
		return profile{rpm: 1500, currentA: 0.8, voltageV: 12.05, tempC: 36, amplitude: 0.15}
	}
}

// Synthesize builds one contract-valid frame and its ground-truth label.
// The vibration math is the Go twin of sentra-core, so simulator output is
// comparable to a frame the firmware would publish.
func Synthesize(regime Regime, seed int64, sequence int, capturedAt time.Time, deviceID, firmware string) (telemetry.Frame, telemetry.Label, error) {
	shape := profileFor(regime)
	rng := xorshift{state: uint64(seed) + uint64(sequence+1)*0x9E3779B97F4A7C15}
	samples := make([]float64, windowSamples)
	for index := range samples {
		if shape.impulsive {
			samples[index] = shape.amplitude * rng.unit()
			if index%19 == 0 {
				samples[index] += 3.2
			}
			continue
		}
		seconds := float64(index) / float64(sampleRateHz)
		frequency := shape.rpm / 60
		samples[index] = shape.amplitude * math.Sin(2*math.Pi*frequency*seconds)
	}
	features, err := dsp.VibrationFeatures(samples, sampleRateHz)
	if err != nil {
		return telemetry.Frame{}, telemetry.Label{}, err
	}
	rssi := -62
	frame := telemetry.Frame{
		SchemaVersion:   telemetry.SchemaVersion,
		DeviceID:        deviceID,
		FirmwareVersion: firmware,
		CapturedAtMs:    capturedAt.UTC().UnixMilli(),
		SampleRateHz:    sampleRateHz,
		WindowSamples:   windowSamples,
		RPM:             shape.rpm,
		CurrentA:        shape.currentA,
		VoltageV:        shape.voltageV,
		TemperatureC:    shape.tempC,
		AccelRMS:        features.AccelRMS,
		AccelPeak:       features.AccelPeak,
		CrestFactor:     features.CrestFactor,
		Kurtosis:        features.Kurtosis,
		DominantFreqHz:  features.DominantFreqHz,
		SensorOK:        true,
		WifiRssiDbm:     &rssi,
	}
	if err := frame.Validate(); err != nil {
		return telemetry.Frame{}, telemetry.Label{}, fmt.Errorf("synthesized frame rejected: %w", err)
	}
	label := telemetry.Label{
		SchemaVersion: telemetry.SchemaVersion,
		DeviceID:      deviceID,
		CapturedAtMs:  frame.CapturedAtMs,
		Condition:     string(regime),
	}
	if err := label.Validate(); err != nil {
		return telemetry.Frame{}, telemetry.Label{}, err
	}
	return frame, label, nil
}

type xorshift struct{ state uint64 }

func (r *xorshift) unit() float64 {
	r.state ^= r.state << 13
	r.state ^= r.state >> 7
	r.state ^= r.state << 17
	if r.state == 0 {
		r.state = 0xA5A5A5A5A5A5A5A5
	}
	return float64(r.state%10_000)/5_000 - 1
}
