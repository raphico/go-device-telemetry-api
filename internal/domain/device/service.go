package device

import (
	"context"

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
	name,
	status,
	deviceType string,
	metadata map[string]any,
) (*Device, error) {
	device, err := NewDevice(userId, name, status, deviceType, metadata)
	if err != nil {
		return nil, err
	}

	if err = s.repo.Create(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

func (s *Service) GetDevice(ctx context.Context, id DeviceID, userId user.UserID) (*Device, error) {
	device, err := s.repo.FindById(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return device, nil
}
