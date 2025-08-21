package telemetry

import (
	"time"

	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
)

type Telemetry struct {
	ID            TelemetryID
	DeviceID      device.DeviceID
	TelemetryType TelemetryType
	Payload       Payload
	RecordedAt    time.Time
}

func NewTelemetry(
	deviceId device.DeviceID,
	telemetryType TelemetryType,
	payload Payload,
) *Telemetry {
	return &Telemetry{
		DeviceID:      deviceId,
		TelemetryType: telemetryType,
		Payload:       payload,
	}
}
