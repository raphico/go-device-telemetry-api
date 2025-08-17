package token

import "context"

type Repository interface {
	Create(ctx context.Context, t *Token) error
}
