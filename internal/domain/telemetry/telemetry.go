package telemetry

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
)

type Telemetry struct {
	ID            TelemetryID
	DeviceID      device.DeviceID
	TelemetryType TelemetryType
	Payload       Payload
	RecordedAt    RecordedAt
	CreatedAt     time.Time
}

func NewTelemetry(
	deviceId device.DeviceID,
	telemetryType TelemetryType,
	payload Payload,
	recordedAt RecordedAt,
) *Telemetry {
	return &Telemetry{
		DeviceID:      deviceId,
		TelemetryType: telemetryType,
		Payload:       payload,
		RecordedAt:    recordedAt,
	}
}

func RehydrateTelemetry(
	id uuid.UUID,
	deviceID uuid.UUID,
	telemetryType string,
	payloadBytes []byte,
	recordedAt time.Time,
	createdAt time.Time,
) (*Telemetry, error) {
	t, err := NewTelemetryType(telemetryType)
	if err != nil {
		return nil, fmt.Errorf("corrupt telemetry type: %w", err)
	}

	var payload map[string]any
	if payloadBytes == nil {
		payload = make(map[string]any)
	} else {
		err := json.Unmarshal(payloadBytes, &payload)
		if err != nil {
			return nil, fmt.Errorf("corrupt payload: %w", err)
		}
	}

	r, err := RecordedAtFromTime(recordedAt)
	if err != nil {
		return nil, fmt.Errorf("corrupt recorded_at: %w", err)
	}

	return &Telemetry{
		ID:            TelemetryID(id),
		DeviceID:      device.DeviceID(deviceID),
		TelemetryType: t,
		Payload:       payload,
		RecordedAt:    r,
		CreatedAt:     createdAt,
	}, nil
}
