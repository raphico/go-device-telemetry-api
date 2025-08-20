package device

import (
	"context"

	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateDevice(
	ctx context.Context,
	userId user.UserID,
	name Name,
	status Status,
	deviceType DeviceType,
	metadata map[string]any,
) (*Device, error) {
	dev := NewDevice(userId, name, status, deviceType, metadata)

	if err := s.repo.Create(ctx, dev); err != nil {
		return nil, err
	}

	return dev, nil
}

func (s *Service) GetDevice(ctx context.Context, id DeviceID, userId user.UserID) (*Device, error) {
	device, err := s.repo.FindById(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *Service) ListUserDevices(
	ctx context.Context,
	userID user.UserID,
	limit int,
	cursor *pagination.Cursor,
) ([]*Device, *pagination.Cursor, error) {
	return s.repo.FindDevices(ctx, userID, limit, cursor)
}
