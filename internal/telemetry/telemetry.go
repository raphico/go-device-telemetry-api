package telemetry

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/raphico/go-device-telemetry-api/internal/device"
)

var (
	telemetryTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// ---------- Types ----------

type TelemetryID uuid.UUID

type Payload map[string]any

type TelemetryType struct {
	value string
}

type RecordedAt struct {
	value time.Time
}

type Telemetry struct {
	ID            TelemetryID
	DeviceID      device.DeviceID
	TelemetryType TelemetryType
	Payload       Payload
	RecordedAt    RecordedAt
	CreatedAt     time.Time
}

// ---------- TelemetryID ----------

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

// ---------- TelemetryType ----------

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

// ---------- Payload ----------

func NewPayload(raw any) (Payload, error) {
	if raw == nil {
		return Payload{}, errors.New("telemetry payload is required")
	}

	p, ok := raw.(map[string]any)
	if !ok {
		return Payload{}, errors.New("telemetry payload must be a valid JSON object")
	}

	if len(p) == 0 {
		return Payload{}, errors.New("telemetry payload cannot be empty")
	}

	return Payload(p), nil
}

// ---------- RecordedAt ----------

func NewRecordedAt(raw string) (RecordedAt, error) {
	now := time.Now().UTC()

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return RecordedAt{}, errors.New("telemetry recorded_at is required")
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil || t.IsZero() {
		return RecordedAt{}, errors.New("invalid recorded_at")
	}

	if t.After(now.Add(5 * time.Minute)) { // reject suspicious future values
		return RecordedAt{}, errors.New("telemetry recorded_at cannot be in the future")
	}

	return RecordedAt{value: t.UTC()}, nil
}

func RecordedAtFromTime(t time.Time) (RecordedAt, error) {
	now := time.Now().UTC()

	if t.IsZero() {
		return RecordedAt{}, errors.New("invalid recorded_at")
	}

	if t.After(now.Add(5 * time.Minute)) { // reject suspicious future values
		return RecordedAt{}, errors.New("telemetry recorded_at cannot be in the future")
	}

	return RecordedAt{value: t.UTC()}, nil
}

func (r RecordedAt) Time() time.Time {
	return r.value
}

// ---------- Telemetry ----------

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

// ---------- Rehydration ----------

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

	var payload Payload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("corrupt payload: %w", err)
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
