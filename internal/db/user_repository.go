package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/user"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		u.Username.String(),
		u.Email.String(),
		u.Password.Hash(),
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return user.ErrEmailAlreadyExists
			case "users_username_key":
				return user.ErrUsernameTaken
			}
		}

		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	var (
		id                   uuid.UUID
		emailStr             string
		usernameStr          string
		passwordHash         []byte
		createdAt, updatedAt time.Time
	)

	query := `
		SELECT id, email, username, password_hash, created_at, updated_at 
		FROM users
		WHERE users.email = $1
	`

	err := r.db.QueryRow(ctx, query, email.String()).Scan(
		&id,
		&emailStr,
		&usernameStr,
		&passwordHash,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user.RehydrateUser(id, emailStr, usernameStr, passwordHash, createdAt, updatedAt)
}
