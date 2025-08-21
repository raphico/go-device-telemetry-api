package device

import "errors"

type Status string

const (
	StatusOffline Status = "offline"
	StatusOnline  Status = "online"
)

func NewStatus(value string) (Status, error) {
	switch value {
	case string(StatusOffline):
		return StatusOffline, nil
	case string(StatusOnline):
		return StatusOnline, nil
	default:
		return "", errors.New("device status must be either 'offline' or 'online'")
	}
}
