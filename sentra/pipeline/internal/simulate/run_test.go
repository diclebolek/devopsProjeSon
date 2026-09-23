package simulate

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

type capture struct {
	frames []telemetry.Frame
	labels []telemetry.Label
	status []telemetry.Status
}

func (c *capture) PublishTelemetry(_ context.Context, frame telemetry.Frame) error {
	c.frames = append(c.frames, frame)
	return nil
}

func (c *capture) PublishLabel(_ context.Context, label telemetry.Label) error {
	c.labels = append(c.labels, label)
	return nil
}

func (c *capture) PublishStatus(_ context.Context, _ string, status telemetry.Status) error {
	c.status = append(c.status, status)
	return nil
}

func TestRunPublishesAClosedCycle(t *testing.T) {
	var recorded capture
	err := Run(context.Background(), Config{
		DeviceID:        "motor-lab-01",
		FirmwareVersion: "0.1.0",
		Regime:          RegimeCycle,
		RegimePeriod:    1,
		Interval:        0,
		MaxFrames:       4,
		Seed:            1,
	}, &recorded, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	if len(recorded.frames) != 4 || len(recorded.labels) != 4 {
		t.Fatalf("frames %d labels %d", len(recorded.frames), len(recorded.labels))
	}
	want := []string{"NORMAL", "OVERLOAD", "UNBALANCED", "HIGH_VIBRATION"}
	for i, condition := range want {
		if recorded.labels[i].Condition != condition {
			t.Fatalf("label %d = %s", i, recorded.labels[i].Condition)
		}
		if err := recorded.frames[i].Validate(); err != nil {
			t.Fatal(err)
		}
	}
	if len(recorded.status) < 2 || !recorded.status[0].Online || recorded.status[len(recorded.status)-1].Online {
		t.Fatalf("status bookends %#v", recorded.status)
	}
	if recorded.frames[0].CapturedAtMs <= 0 {
		t.Fatal("missing timestamp")
	}
	_ = time.Second
}
