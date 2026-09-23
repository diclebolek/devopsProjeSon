package store

import (
	"context"
	"time"

	"github.com/diclebolek/devopsProjeSon/sentra/pipeline/internal/telemetry"
)

// FreshAfter is how recently a node must have reported to count as live.
const FreshAfter = 15 * time.Second

// Store is the persistence boundary. Memory is the lab default; Postgres is production.
type Store interface {
	InsertTelemetry(ctx context.Context, frame telemetry.Frame, ingestedAt time.Time) error
	UpsertLabel(ctx context.Context, label telemetry.Label) error
	UpsertStatus(ctx context.Context, deviceID string, status telemetry.Status, seenAt time.Time) error
	ListTelemetry(ctx context.Context, deviceID string, limit int, now time.Time) ([]telemetry.Reading, error)
	ListDevices(ctx context.Context, now time.Time) ([]telemetry.Device, error)
	Ping(ctx context.Context) error
}
