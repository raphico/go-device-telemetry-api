package device

import (
	"time"

	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
)

type Device struct {
	ID         DeviceID
	UserID     user.UserID
	Name       Name
	Status     Status
	DeviceType DeviceType
	Metadata   map[string]any
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewDevice(userId user.UserID, name, status, deviceType string, metadata map[string]any) (*Device, error) {
	n, err := NewName(name)
	if err != nil {
		return nil, err
	}

	s, err := NewStatus(status)
	if err != nil {
		return nil, err
	}

	dt, err := NewDeviceType(deviceType)
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = make(map[string]any)
	}

	return &Device{
		UserID:     userId,
		Name:       n,
		Status:     s,
		DeviceType: dt,
		Metadata:   metadata,
	}, nil
}
