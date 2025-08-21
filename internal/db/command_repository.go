package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/domain/command"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
)

type CommandRepository struct {
	db *pgxpool.Pool
}

func NewCommandRepository(db *pgxpool.Pool) *CommandRepository {
	return &CommandRepository{db: db}
}

func (r *CommandRepository) Create(ctx context.Context, c *command.Command) error {
	jsonPayload, err := json.Marshal(c.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO commands (device_id, command_name, payload)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, status
	`

	var status string
	err = r.db.QueryRow(
		ctx,
		query,
		c.DeviceID,
		c.Name.String(),
		jsonPayload,
	).Scan(&c.ID, &c.CreatedAt, &status)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == "23503" && pgError.ConstraintName == "commands_device_id_fkey" {
				return device.ErrDeviceNotFound
			}
		}
		return fmt.Errorf("failed to insert command: %w", err)
	}

	if err := c.UpdateStatus(status); err != nil {
		return fmt.Errorf("corrupt status: %w", err)
	}

	return nil
}
