package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (r *TokenRepository) FindValidTokenByHash(ctx context.Context, hash []byte, scope string) (*token.Token, error) {
	t := &token.Token{}
	query := `
		SELECT id, token_hash, user_id, scope, revoked, expires_at, last_used_at, created_at
		FROM tokens
		WHERE token_hash = $1 
			AND scope = $2
			AND revoked = false 
			AND expires_at > now()
	`

	err := r.db.QueryRow(ctx, query, hash, scope).Scan(
		&t.ID,
		&t.Hash,
		&t.UserID,
		&t.Scope,
		&t.Revoked,
		&t.ExpiresAt,
		&t.LastUsedAt,
		&t.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, token.ErrTokenNotFound
		}

		return nil, fmt.Errorf("failed to find token: %w", err)
	}

	return t, nil
}

func (r *TokenRepository) Revoke(ctx context.Context, scope string, hash []byte) error {
	query := `
		UPDATE tokens
		SET revoked = true
		WHERE token_hash = $1
			AND scope = $2
			AND revoked = false
	`

	tag, err := r.db.Exec(ctx, query, hash, scope)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return token.ErrTokenNotFound
	}

	return nil
}

func (r *TokenRepository) UpdateLastUsed(ctx context.Context, hash []byte) error {
	query := `UPDATE tokens SET last_used_at = now() WHERE token_hash = $1 AND revoked = false`
	tag, err := r.db.Exec(ctx, query, hash)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return token.ErrTokenNotFound
	}

	return nil
}
