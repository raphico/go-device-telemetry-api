package command

import (
	"context"

	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
)

type Repository interface {
	Create(ctx context.Context, c *Command) error
	FindCommands(
		ctx context.Context,
		deviceID device.DeviceID,
		limit int,
		cursor *pagination.Cursor,
	) ([]*Command, *pagination.Cursor, error)
}
