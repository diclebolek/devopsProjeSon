package telemetry

import (
	"os"
	"testing"
)

func TestGoldenFrameValidates(t *testing.T) {
	raw, err := os.ReadFile("../../../testdata/frame_normal.json")
	if err != nil {
		t.Fatal(err)
	}
	frame, err := DecodeFrame(raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := frame.Validate(); err != nil {
		t.Fatal(err)
	}
	if frame.DeviceID != "motor-lab-01" {
		t.Fatalf("device id %s", frame.DeviceID)
	}
	if frame.WifiRssiDbm == nil || *frame.WifiRssiDbm != -61 {
		t.Fatalf("rssi %#v", frame.WifiRssiDbm)
	}
}

func TestRejectsNonFiniteAndBadIdentity(t *testing.T) {
	frame := validFrame()
	frame.CurrentA = -1
	if err := frame.Validate(); err == nil {
		t.Fatal("expected current rejection")
	}
	frame = validFrame()
	frame.Kurtosis = 1e9
	if err := frame.Validate(); err == nil {
		t.Fatal("expected kurtosis rejection")
	}
	frame = validFrame()
	frame.DeviceID = "Motor"
	if err := frame.Validate(); err == nil {
		t.Fatal("expected device id rejection")
	}
}

func TestTopicParser(t *testing.T) {
	id, kind, err := ParseTopic("sentra/v1/motor-lab-01/telemetry")
	if err != nil || id != "motor-lab-01" || kind != "telemetry" {
		t.Fatalf("got %s %s %v", id, kind, err)
	}
	for _, topic := range []string{
		"sentra/v2/motor-lab-01/telemetry",
		"sentra/v1/telemetry",
		"sentra/v1/Motor/telemetry",
		"sentra/v1/motor-lab-01/waveform",
	} {
		if _, _, err := ParseTopic(topic); err == nil {
			t.Fatalf("accepted %s", topic)
		}
	}
	if err := SameDevice("motor-lab-01", "motor-lab-02"); err == nil {
		t.Fatal("expected mismatch")
	}
}

func TestLabelVocabulary(t *testing.T) {
	label := Label{SchemaVersion: 1, DeviceID: "motor-lab-01", CapturedAtMs: 10, Condition: "OVERLOAD"}
	if err := label.Validate(); err != nil {
		t.Fatal(err)
	}
	label.Condition = "BEARING"
	if err := label.Validate(); err == nil {
		t.Fatal("expected closed vocabulary")
	}
}

func validFrame() Frame {
	rssi := -61
	return Frame{
		SchemaVersion:   1,
		DeviceID:        "motor-lab-01",
		FirmwareVersion: "0.1.0",
		CapturedAtMs:    1_710_000_000_000,
		SampleRateHz:    1000,
		WindowSamples:   256,
		RPM:             1500,
		CurrentA:        0.8,
		VoltageV:        12,
		TemperatureC:    35.5,
		AccelRMS:        0.12,
		AccelPeak:       0.2,
		CrestFactor:     1.6,
		Kurtosis:        1.5,
		DominantFreqHz:  25,
		SensorOK:        true,
		WifiRssiDbm:     &rssi,
	}
}
