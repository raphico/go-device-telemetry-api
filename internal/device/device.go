package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/raphico/go-device-telemetry-api/internal/user"
)

var (
	StatusOffline = Status{"offline"}
	StatusOnline  = Status{"online"}

	deviceTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	nameRegex       = regexp.MustCompile(`^[a-zA-Z0-9 _.-]+$`)
)

// ---------- Types ----------

type DeviceID uuid.UUID

type Name struct {
	value string
}

type Status struct {
	value string
}

type DeviceType struct {
	value string
}

type Metadata map[string]any

type Device struct {
	ID         DeviceID
	UserID     user.UserID
	Name       Name
	Status     Status
	DeviceType DeviceType
	Metadata   map[string]any
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ---------- DeviceID ----------

func NewDeviceID(id string) (DeviceID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return DeviceID(uuid.Nil), err
	}

	return DeviceID(parsed), nil
}

func (u DeviceID) String() string {
	return uuid.UUID(u).String()
}

// ---------- Name ----------

func NewName(value string) (Name, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Name{}, errors.New("device name is required")
	}

	if len(value) < 3 {
		return Name{}, errors.New("device name must be at least 3 characters")
	}
	if len(value) > 50 {
		return Name{}, errors.New("device name must be at most 50 characters")
	}

	if !nameRegex.MatchString(value) {
		return Name{}, errors.New("device name may only contain letters, numbers, underscores, periods, or hyphens")
	}

	return Name{value: value}, nil
}

func (n Name) String() string {
	return n.value
}

// ---------- Metadata ----------

func NewMetadata(raw any) (Metadata, error) {
	if raw == nil {
		return Metadata{}, errors.New("device metadata is required")
	}

	p, ok := raw.(map[string]any)
	if !ok {
		return Metadata{}, errors.New("device metadata must be a valid JSON object")
	}

	if len(p) == 0 {
		return Metadata{}, errors.New("device metadata cannot be empty")
	}

	return Metadata(p), nil
}

// ---------- Status ----------

func NewStatus(value string) (Status, error) {
	switch value {
	case StatusOffline.value:
		return StatusOffline, nil
	case StatusOnline.value:
		return StatusOnline, nil
	default:
		return Status{}, fmt.Errorf("invalid device status: %s", value)
	}
}

func (s Status) String() string {
	return s.value
}

// ---------- DeviceType ----------

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

// ---------- Device ----------

func NewDevice(
	userId user.UserID,
	name Name,
	status Status,
	deviceType DeviceType,
	metadata map[string]any,
) *Device {
	if metadata == nil {
		metadata = make(map[string]any)
	}

	return &Device{
		UserID:     userId,
		Name:       name,
		Status:     status,
		DeviceType: deviceType,
		Metadata:   metadata,
	}
}

func (d *Device) UpdateName(n Name) {
	d.Name = n
}

func (d *Device) UpdateDeviceType(dt DeviceType) {
	d.DeviceType = dt
}

func (d *Device) UpdateMetadata(m map[string]any) {
	d.Metadata = m
}

// ---------- Rehydration ----------

func RehydrateDevice(
	id uuid.UUID,
	userID uuid.UUID,
	name string,
	deviceType string,
	status string,
	metadataBytes []byte,
	createdAt time.Time,
	updatedAt time.Time,
) (*Device, error) {
	n, err := NewName(name)
	if err != nil {
		return nil, fmt.Errorf("corrupt device name: %w", err)
	}

	s, err := NewStatus(status)
	if err != nil {
		return nil, fmt.Errorf("corrupt device status: %w", err)
	}

	dt, err := NewDeviceType(deviceType)
	if err != nil {
		return nil, fmt.Errorf("corrupt device type: %w", err)
	}

	var metadata map[string]any
	if metadataBytes == nil {
		metadata = make(map[string]any)
	} else {
		err := json.Unmarshal(metadataBytes, &metadata)
		if err != nil {
			return nil, fmt.Errorf("corrupt metadata: %w", err)
		}
	}

	return &Device{
		ID:         DeviceID(id),
		UserID:     user.UserID(userID),
		Name:       n,
		Status:     s,
		DeviceType: dt,
		Metadata:   metadata,
		UpdatedAt:  updatedAt,
		CreatedAt:  createdAt,
	}, nil
}
