package telemetry

import (
	"errors"
	"regexp"
	"strings"
)

type TelemetryType struct {
	value string
}

var telemetryTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func NewTelemetryType(t string) (TelemetryType, error) {
	t = strings.TrimSpace(t)

	if len(t) == 0 {
		return TelemetryType{}, errors.New("telemetry type is required")
	}

	if len(t) < 3 {
		return TelemetryType{}, errors.New("telemetry type must be at least 3 characters")
	}

	if len(t) > 50 {
		return TelemetryType{}, errors.New("telemetry type must be at most 50 characters")
	}

	if matches := telemetryTypeRegex.MatchString(t); !matches {
		return TelemetryType{}, errors.New("telemetry type can only contain letters, numbers, _ and -")
	}

	return TelemetryType{value: t}, nil
}

func (t TelemetryType) String() string {
	return t.value
}
