package token

import "errors"

var (
	ErrTokenGenerationFailed = errors.New("failed to generate token")
	ErrTokenNotFound         = errors.New("token not found")
	ErrTokenAlreadyExists    = errors.New("token already exists")
	ErrInvalidToken          = errors.New("invalid token")
	ErrExpiredToken          = errors.New("token expired")
	ErrWrongTokenType        = errors.New("wrong signing method")
)
