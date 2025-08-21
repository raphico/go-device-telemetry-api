package telemetry

import "github.com/google/uuid"

type TelemetryID uuid.UUID

func NewTelemetryID(id string) (TelemetryID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return TelemetryID(uuid.Nil), err
	}

	return TelemetryID(parsed), nil
}

func (t TelemetryID) String() string {
	return uuid.UUID(t).String()
}
