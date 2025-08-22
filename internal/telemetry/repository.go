package telemetry

import (
	"context"

	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/device"
)

type Repository interface {
	Create(ctx context.Context, t *Telemetry) error
	FindTelemetry(
		ctx context.Context,
		deviceID device.DeviceID,
		limit int,
		cursor *pagination.Cursor,
	) ([]*Telemetry, *pagination.Cursor, error)
}
