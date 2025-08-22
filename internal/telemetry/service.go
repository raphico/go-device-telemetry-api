package telemetry

import (
	"context"

	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/device"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTelemetry(
	ctx context.Context,
	deviceID device.DeviceID,
	telemetryType TelemetryType,
	payload Payload,
	recordedAt RecordedAt,
) (*Telemetry, error) {
	telemetry := NewTelemetry(deviceID, telemetryType, payload, recordedAt)

	err := s.repo.Create(ctx, telemetry)
	if err != nil {
		return nil, err
	}

	return telemetry, nil
}

func (s *Service) ListDeviceTelemetry(
	ctx context.Context,
	deviceID device.DeviceID,
	limit int,
	cursor *pagination.Cursor,
) ([]*Telemetry, *pagination.Cursor, error) {
	return s.repo.FindTelemetry(ctx, deviceID, limit, cursor)
}
