package command

import (
	"errors"
	"strings"
	"time"
)

type ExecutedAt struct {
	value time.Time
	valid bool
}

func NewExecutedAt(raw string) (ExecutedAt, error) {
	now := time.Now().UTC()

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ExecutedAt{}, errors.New("command executed_at is required")
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil || t.IsZero() {
		return ExecutedAt{}, errors.New("invalid executed_at")
	}

	if t.After(now.Add(5 * time.Minute)) { // reject suspicious future values
		return ExecutedAt{}, errors.New("command executed_at cannot be in the future")
	}

	return ExecutedAt{value: t.UTC(), valid: true}, nil
}

func ExecutedAtFromTime(t time.Time) (ExecutedAt, error) {
	if t.IsZero() {
		return ExecutedAt{valid: false}, nil // NULL
	}

	now := time.Now().UTC()
	if t.After(now.Add(5 * time.Minute)) {
		return ExecutedAt{}, errors.New("command executed_at cannot be in the future")
	}

	return ExecutedAt{value: t.UTC(), valid: true}, nil
}

func (e ExecutedAt) Time() time.Time {
	return e.value
}

func (e ExecutedAt) Valid() bool {
	return e.valid
}
