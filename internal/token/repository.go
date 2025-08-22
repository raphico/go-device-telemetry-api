package token

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, t *Token) error
	FindValidTokenByHash(ctx context.Context, hash []byte, scope string) (*Token, error)
	Revoke(ctx context.Context, scope string, hash []byte) error
	UpdateLastUsed(ctx context.Context, hash []byte) error
}
