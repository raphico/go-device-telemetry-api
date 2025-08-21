package command

import (
	"context"

	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCommand(
	ctx context.Context,
	deviceID device.DeviceID,
	name Name,
	payload Payload,
) (*Command, error) {
	command := NewCommand(deviceID, name, payload)

	err := s.repo.Create(ctx, command)
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (s *Service) ListDeviceCommands(
	ctx context.Context,
	deviceID device.DeviceID,
	limit int,
	cursor *pagination.Cursor,
) ([]*Command, *pagination.Cursor, error) {
	return s.repo.FindCommands(ctx, deviceID, limit, cursor)
}
