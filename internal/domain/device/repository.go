package device

import (
	"context"

	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, device *Device) error
	FindById(ctx context.Context, id DeviceID, userId user.UserID) (*Device, error)
}
