package telemetry

import (
	"errors"
	"strings"
	"time"
)

type RecordedAt struct {
	value time.Time
}

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

func (r RecordedAt) Time() time.Time {
	return r.value
}
