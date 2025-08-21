package device

import (
	"context"

	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, device *Device) error
	FindById(ctx context.Context, id DeviceID, userId user.UserID) (*Device, error)
	FindDevices(ctx context.Context, userId user.UserID, limit int, cursor *pagination.Cursor) ([]*Device, *pagination.Cursor, error)
	Update(ctx context.Context, dev *Device) error
}
