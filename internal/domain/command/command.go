package command

import (
	"time"

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
