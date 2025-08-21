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
	cmd := NewCommand(deviceID, name, payload)

	err := s.repo.Create(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (s *Service) ListDeviceCommands(
	ctx context.Context,
	deviceID device.DeviceID,
	limit int,
	cursor *pagination.Cursor,
) ([]*Command, *pagination.Cursor, error) {
	return s.repo.FindCommands(ctx, deviceID, limit, cursor)
}

func (s *Service) UpdateCommandStatus(
	ctx context.Context, 
	id CommandID,
	deviceID device.DeviceID,
	status Status,
	executedAt ExecutedAt,	
) (*Command, error) {
	cmd, err := s.repo.FindById(ctx, id, deviceID)
	if err != nil {
		return nil, err
	}

	cmd.UpdateStatus(status)
	cmd.UpdateExecutedAt(executedAt)

	err = s.repo.UpdateStatus(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
