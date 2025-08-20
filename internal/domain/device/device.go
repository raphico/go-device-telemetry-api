package device

import (
	"encoding/json"
	"fmt"
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

func NewDevice(
	userId user.UserID,
	name Name,
	status Status,
	deviceType DeviceType,
	metadata map[string]any,
) *Device {
	if metadata == nil {
		metadata = make(map[string]any)
	}

	return &Device{
		UserID:     userId,
		Name:       name,
		Status:     status,
		DeviceType: deviceType,
		Metadata:   metadata,
	}
}

func RehydrateDevice(
	id DeviceID,
	userID user.UserID,
	name string,
	deviceType string,
	status string,
	metadataBytes []byte,
	createdAt time.Time,
	updatedAt time.Time,
) (*Device, error) {
	n, err := NewName(name)
	if err != nil {
		return nil, fmt.Errorf("corrupt device name: %w", err)
	}

	s, err := NewStatus(status)
	if err != nil {
		return nil, fmt.Errorf("corrupt device status: %w", err)
	}

	dt, err := NewDeviceType(deviceType)
	if err != nil {
		return nil, fmt.Errorf("corrupt device type: %w", err)
	}

	var metadata map[string]any
	if metadataBytes == nil {
		metadata = make(map[string]any)
	} else {
		err := json.Unmarshal(metadataBytes, &metadata)
		if err != nil {
			return nil, fmt.Errorf("corrupt metadata: %w", err)
		}
	}

	return &Device{
		ID:         id,
		UserID:     userID,
		Name:       n,
		Status:     s,
		DeviceType: dt,
		Metadata:   metadata,
		UpdatedAt:  updatedAt,
		CreatedAt:  createdAt,
	}, nil
}
