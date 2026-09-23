package dsp

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

func TestSharedOracle(t *testing.T) {
	raw, err := os.ReadFile("../../../testdata/dsp_case.json")
	if err != nil {
		t.Fatal(err)
	}
	var file struct {
		Sine struct {
			SampleRateHz           int     `json:"sample_rate_hz"`
			FrequencyHz            float64 `json:"frequency_hz"`
			Amplitude              float64 `json:"amplitude"`
			Offset                 float64 `json:"offset"`
			Count                  int     `json:"count"`
			ExpectedRMS            float64 `json:"expected_rms"`
			ExpectedPeak           float64 `json:"expected_peak"`
			ExpectedCrestFactor    float64 `json:"expected_crest_factor"`
			ExpectedKurtosis       float64 `json:"expected_kurtosis"`
			ExpectedDominantFreqHz float64 `json:"expected_dominant_freq_hz"`
		} `json:"sine"`
		Fixed struct {
			SampleRateHz           int       `json:"sample_rate_hz"`
			Samples                []float64 `json:"samples"`
			ExpectedRMS            float64   `json:"expected_rms"`
			ExpectedPeak           float64   `json:"expected_peak"`
			ExpectedCrestFactor    float64   `json:"expected_crest_factor"`
			ExpectedKurtosis       float64   `json:"expected_kurtosis"`
			ExpectedDominantFreqHz float64   `json:"expected_dominant_freq_hz"`
		} `json:"fixed"`
	}
	if err := json.Unmarshal(raw, &file); err != nil {
		t.Fatal(err)
	}

	sine := make([]float64, file.Sine.Count)
	for i := range sine {
		phase := 2 * math.Pi * file.Sine.FrequencyHz * float64(i) / float64(file.Sine.SampleRateHz)
		sine[i] = file.Sine.Offset + file.Sine.Amplitude*math.Sin(phase)
	}
	features, err := VibrationFeatures(sine, file.Sine.SampleRateHz)
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, features.AccelRMS, file.Sine.ExpectedRMS)
	assertClose(t, features.AccelPeak, file.Sine.ExpectedPeak)
	assertClose(t, features.CrestFactor, file.Sine.ExpectedCrestFactor)
	assertClose(t, features.Kurtosis, file.Sine.ExpectedKurtosis)
	assertClose(t, features.DominantFreqHz, file.Sine.ExpectedDominantFreqHz)

	fixed, err := VibrationFeatures(file.Fixed.Samples, file.Fixed.SampleRateHz)
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, fixed.AccelRMS, file.Fixed.ExpectedRMS)
	assertClose(t, fixed.AccelPeak, file.Fixed.ExpectedPeak)
	assertClose(t, fixed.CrestFactor, file.Fixed.ExpectedCrestFactor)
	assertClose(t, fixed.Kurtosis, file.Fixed.ExpectedKurtosis)
	assertClose(t, fixed.DominantFreqHz, file.Fixed.ExpectedDominantFreqHz)
}

func assertClose(t *testing.T, actual, expected float64) {
	t.Helper()
	if math.Abs(actual-expected) > 1e-9 {
		t.Fatalf("actual %v expected %v", actual, expected)
	}
}
