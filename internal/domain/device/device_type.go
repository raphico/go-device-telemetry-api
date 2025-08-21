package device

import (
	"errors"
	"regexp"
	"strings"
)

type DeviceType struct {
	value string
}

var deviceTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func NewDeviceType(t string) (DeviceType, error) {
	t = strings.TrimSpace(t)

	if len(t) == 0 {
		return DeviceType{}, errors.New("device type is required")
	}

	if len(t) < 3 {
		return DeviceType{}, errors.New("device type must be at least 3 characters")
	}

	if len(t) > 50 {
		return DeviceType{}, errors.New("device type must be at most 50 characters")
	}

	if matches := deviceTypeRegex.MatchString(t); !matches {
		return DeviceType{}, errors.New("device type can only contain letters, numbers, _ and -")
	}

	return DeviceType{value: t}, nil
}

func (t DeviceType) String() string {
	return t.value
}
