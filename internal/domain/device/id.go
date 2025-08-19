package device

import "github.com/google/uuid"

type DeviceID uuid.UUID

func NewDeviceID(id string) (DeviceID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return DeviceID(uuid.Nil), err
	}

	return DeviceID(parsed), nil
}

func (u DeviceID) String() string {
	return uuid.UUID(u).String()
}
