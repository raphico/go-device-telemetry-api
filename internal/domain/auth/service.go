package auth

import (
	"context"

	"github.com/raphico/go-device-telemetry-api/internal/domain/token"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
)

type Service struct {
	user  *user.Service
	token *token.Service
}

func NewService(userService *user.Service, tokenService *token.Service) *Service {
	return &Service{
		user:  userService,
		token: tokenService,
	}
}

func (s *Service) Register(
	ctx context.Context,
	username user.Username,
	email user.Email,
	password user.Password,
) (*user.User, error) {
	return s.user.RegisterUser(ctx, username, email, password)
}

func (s *Service) Login(ctx context.Context, email user.Email, rawPassword string) (string, *token.Token, error) {
	u, err := s.user.AuthenticateUser(ctx, email, rawPassword)
	if err != nil {
		return "", nil, err
	}

	accessToken, err := s.token.GenerateAccessToken(u.ID)
	if err != nil {
		return "", nil, err
	}

	refreshToken, err := s.token.CreateRefreshToken(ctx, u.ID)
	if err != nil {
		return "", nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Refresh(ctx context.Context, refreshTok string) (string, *token.Token, error) {
	return s.token.RotateTokens(ctx, refreshTok)
}

func (s *Service) Logout(ctx context.Context, refreshTok string) error {
	return s.token.RevokeRefreshToken(ctx, refreshTok)
}
