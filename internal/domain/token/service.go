package token

import (
	"context"
	"fmt"
	"time"

	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
)

type Service struct {
	repo            Repository
	jwtGen          JWTGenerator
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewService(
	jwtGen JWTGenerator,
	repo Repository,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Service {
	return &Service{
		repo:            repo,
		jwtGen:          jwtGen,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *Service) GenerateAccessToken(userId user.UserID) (string, error) {
	token, err := s.jwtGen.Generate(userId.String(), s.accessTokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) ValidateAccessToken(tokenStr string) (user.UserID, error) {
	claims, err := s.jwtGen.Validate(tokenStr)
	if err != nil {
		return user.UserID{}, err
	}

	UserId, err := user.NewUserID(claims.UserID)
	if err != nil {
		return user.UserID{}, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return UserId, nil
}

func (s *Service) CreateRefreshToken(ctx context.Context, userId user.UserID) (*Token, error) {
	token, err := NewToken(userId, s.refreshTokenTTL, "auth")
	if err != nil {
		return nil, err
	}

	if err = s.repo.Create(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) RotateTokens(ctx context.Context, refreshTok string) (string, *Token, error) {
	hash := HashPlaintext(refreshTok)

	tokenRecord, err := s.repo.FindValidTokenByHash(ctx, hash, "auth")
	if err != nil {
		fmt.Println(err)
		return "", nil, err
	}

	if err := s.repo.UpdateLastUsed(ctx, hash); err != nil {
		return "", nil, err
	}

	if err := s.repo.Revoke(ctx, "auth", tokenRecord.Hash); err != nil {
		return "", nil, err
	}

	accessToken, err := s.GenerateAccessToken(tokenRecord.UserID)
	if err != nil {
		return "", nil, err
	}

	token, err := s.CreateRefreshToken(ctx, tokenRecord.UserID)
	if err != nil {
		return "", nil, err
	}

	return accessToken, token, nil
}
