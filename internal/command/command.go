package command

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
	StatusPending  = Status{"pending"}
	StatusExecuted = Status{"executed"}
	StatusFailed   = Status{"failed"}

	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9 _.-]+$`)
)

// ---------- Types ----------

type CommandID uuid.UUID

type Payload map[string]any

type Name struct {
	value string
}

type Status struct {
	value string
}

type ExecutedAt struct {
	value time.Time
	valid bool
}

type Command struct {
	ID         CommandID
	DeviceID   device.DeviceID
	Name       Name
	Payload    Payload
	Status     Status
	ExecutedAt ExecutedAt
	CreatedAt  time.Time
}

// ---------- CommandID ----------

func NewCommandID(id string) (CommandID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return CommandID(uuid.Nil), err
	}

	return CommandID(parsed), nil
}

func (c CommandID) String() string {
	return uuid.UUID(c).String()
}

// ---------- Name ----------

func NewName(value string) (Name, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Name{}, errors.New("command name is required")
	}

	if len(value) < 3 {
		return Name{}, errors.New("command name must be at least 3 characters")
	}
	if len(value) > 50 {
		return Name{}, errors.New("command name must be at most 50 characters")
	}

	if !nameRegex.MatchString(value) {
		return Name{}, errors.New("command name may only contain letters, numbers, underscores, periods, or hyphens")
	}

	return Name{value: value}, nil
}

func (n Name) String() string {
	return n.value
}

// ---------- Payload ----------

func NewPayload(raw any) (Payload, error) {
	if raw == nil {
		return Payload{}, errors.New("command payload is required")
	}

	p, ok := raw.(map[string]any)
	if !ok {
		return Payload{}, errors.New("command payload must be a valid JSON object")
	}

	if len(p) == 0 {
		return Payload{}, errors.New("command payload cannot be empty")
	}

	return Payload(p), nil
}

// ---------- Status ----------

func NewStatus(value string) (Status, error) {
	switch value {
	case StatusPending.value:
		return StatusPending, nil
	case StatusExecuted.value:
		return StatusExecuted, nil
	case StatusFailed.value:
		return StatusFailed, nil
	default:
		return Status{}, fmt.Errorf("invalid command status: %s", value)
	}
}

func (s Status) String() string {
	return s.value
}

func (s *Status) SetStatus(value string) error {
	status, err := NewStatus(value)
	if err != nil {
		return err
	}
	*s = status
	return nil
}

// ---------- ExecutedAt ----------

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

// ---------- Command ----------

func NewCommand(
	deviceID device.DeviceID,
	name Name,
	payload Payload,
) *Command {
	return &Command{
		DeviceID: deviceID,
		Name:     name,
		Payload:  payload,
	}
}

func (c *Command) UpdateStatus(status Status) {
	c.Status = status
}

func (c *Command) UpdateExecutedAt(executedAt ExecutedAt) {
	c.ExecutedAt = executedAt
}

// ---------- Rehydration ----------

func RehydrateCommand(
	id uuid.UUID,
	deviceID uuid.UUID,
	name string,
	payloadBytes []byte,
	status string,
	executedAt *time.Time,
	createdAt time.Time,
) (*Command, error) {
	n, err := NewName(name)
	if err != nil {
		return nil, fmt.Errorf("corrupt command name: %w", err)
	}

	var payload Payload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("corrupt payload: %w", err)
	}

	s, err := NewStatus(status)
	if err != nil {
		return nil, fmt.Errorf("corrupt status: %w", err)
	}

	var execAt ExecutedAt
	if executedAt != nil {
		e, err := ExecutedAtFromTime(*executedAt)
		if err != nil {
			return nil, fmt.Errorf("corrupt executed_at: %w", err)
		}
		execAt = e
	} else {
		execAt = ExecutedAt{valid: false}
	}

	return &Command{
		ID:         CommandID(id),
		DeviceID:   device.DeviceID(deviceID),
		Name:       n,
		Payload:    payload,
		Status:     s,
		ExecutedAt: execAt,
		CreatedAt:  createdAt,
	}, nil
}
