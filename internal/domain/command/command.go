package command

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
)

type Command struct {
	ID         CommandID
	DeviceID   device.DeviceID
	Name       Name
	Payload    Payload
	Status     Status
	ExecutedAt *time.Time
	CreatedAt  time.Time
}

func NewCommand(
	deviceID device.DeviceID,
	name Name,
	payload Payload,
) *Command {
	return &Command{
		DeviceID: deviceID,
		Name:     name,
		Payload:  payload,
	}
}

func (c *Command) UpdateStatus(value string) error {
	status, err := NewStatus(value)
	if err != nil {
		return err
	}
	c.Status = status
	return nil
}

func RehydrateCommand(
	id uuid.UUID,
	deviceID uuid.UUID,
	name string,
	payloadBytes []byte,
	status string,
	executedAt *time.Time,
	createdAt time.Time,
) (*Command, error) {
	n, err := NewName(name)
	if err != nil {
		return nil, fmt.Errorf("corrupt command name: %w", err)
	}

	var payload Payload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("corrupt payload: %w", err)
	}

	s, err := NewStatus(status)
	if err != nil {
		return nil, fmt.Errorf("corrupt status: %w", err)
	}

	return &Command{
		ID:         CommandID(id),
		DeviceID:   device.DeviceID(deviceID),
		Name:       n,
		Payload:    payload,
		Status:     s,
		ExecutedAt: executedAt,
		CreatedAt:  createdAt,
	}, nil
}
