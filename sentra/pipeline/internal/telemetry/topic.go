package telemetry

import "strings"

// ParseTopic splits sentra/v1/{device_id}/{telemetry|label|status}.
func ParseTopic(topic string) (deviceID, kind string, err error) {
	parts := strings.Split(topic, "/")
	if len(parts) != 4 || parts[0] != "sentra" || parts[1] != "v1" {
		return "", "", invalid("topic", "expected sentra/v1/{device_id}/{telemetry|label|status}")
	}
	switch parts[3] {
	case "telemetry", "label", "status":
	default:
		return "", "", invalid("topic", "unknown message kind")
	}
	if !deviceIDPattern.MatchString(parts[2]) {
		return "", "", invalid("device_id", "topic device id is not a valid slug")
	}
	return parts[2], parts[3], nil
}

// SameDevice reports whether a payload id matches the topic id.
func SameDevice(topicDeviceID, payloadDeviceID string) error {
	if topicDeviceID != payloadDeviceID {
		return invalid("device_id", "payload device_id does not match the topic")
	}
	return nil
}
