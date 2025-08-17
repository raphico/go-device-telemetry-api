package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/domain/token"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, t *token.Token) error {
	query := `
		INSERT INTO tokens(user_id, token_hash, scope, expires_at)	
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		t.UserID,
		t.Hash,
		t.Scope,
		t.ExpiresAt,
	).Scan(&t.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "tokens_token_hash_key" {
					return token.ErrTokenAlreadyExists
				}
			case "23503": // foreign_key_constraint
				if pgErr.ConstraintName == "tokens_user_id_fkey" {
					return token.ErrUserNotFound
				}
			}
		}

		return fmt.Errorf("failed to insert token: %w", err)
	}

	return nil
}
