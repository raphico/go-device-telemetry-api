package device

import "context"

type Repository interface {
	Create(ctx context.Context, device *Device) error
}
