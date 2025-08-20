package user

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) RegisterUser(
	ctx context.Context,
	username Username,
	email Email,
	password Password,
) (*User, error) {
	user := NewUser(email, username, password)

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) AuthenticateUser(ctx context.Context, email Email, rawPassword string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if matches := user.Password.Matches(rawPassword); !matches {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
