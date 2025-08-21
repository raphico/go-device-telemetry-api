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
	ExecutedAt ExecutedAt
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

func (c *Command) UpdateStatus(status Status) {
	c.Status = status
}

func (c *Command) UpdateExecutedAt(executedAt ExecutedAt) {
	c.ExecutedAt = executedAt
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

	var execAt ExecutedAt
	if executedAt != nil {
		e, err := ExecutedAtFromTime(*executedAt)
		if err != nil {
			return nil, fmt.Errorf("corrupt executed_at: %w", err)
		}
		execAt = e
	} else {
		execAt = ExecutedAt{valid: false}
	}

	return &Command{
		ID:         CommandID(id),
		DeviceID:   device.DeviceID(deviceID),
		Name:       n,
		Payload:    payload,
		Status:     s,
		ExecutedAt: execAt,
		CreatedAt:  createdAt,
	}, nil
}
