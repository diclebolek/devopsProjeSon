package telemetry

import "encoding/json"

func DecodeFrame(payload []byte) (Frame, error) {
	var frame Frame
	if err := json.Unmarshal(payload, &frame); err != nil {
		return Frame{}, invalid("payload", "telemetry is not valid json")
	}
	return frame, nil
}

func DecodeLabel(payload []byte) (Label, error) {
	var label Label
	if err := json.Unmarshal(payload, &label); err != nil {
		return Label{}, invalid("payload", "label is not valid json")
	}
	return label, nil
}

func DecodeStatus(payload []byte) (Status, error) {
	var status Status
	if err := json.Unmarshal(payload, &status); err != nil {
		return Status{}, invalid("payload", "status is not valid json")
	}
	return status, nil
}

func DecodeStatusReport(payload []byte) (StatusReport, error) {
	var report StatusReport
	if err := json.Unmarshal(payload, &report); err != nil {
		return StatusReport{}, invalid("payload", "status is not valid json")
	}
	return report, nil
}
