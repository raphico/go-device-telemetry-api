package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
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

func (r *CommandRepository) FindCommands(
	ctx context.Context,
	deviceID device.DeviceID,
	limit int,
	cursor *pagination.Cursor,
) ([]*command.Command, *pagination.Cursor, error) {
	var query string
	var args []any

	if cursor == nil {
		query = `
    		SELECT id, device_id, command_name, payload, status, executed_at, created_at
    		FROM commands
    		WHERE device_id = $1
    		ORDER BY created_at ASC, id ASC
    		LIMIT $2
    	`
		args = []any{deviceID, limit + 1}
	} else {
		query = `
    		SELECT id, device_id, command_name, payload, status, executed_at, created_at
    		FROM commands
    		WHERE device_id = $1
				AND (created_at, id) > ($2, $3)
    		ORDER BY created_at ASC, id ASC
    		LIMIT $4
		`
		args = []any{deviceID, cursor.CreatedAt, cursor.ID, limit + 1}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query commands: %w", err)
	}

	var result []*command.Command
	for rows.Next() {
		var (
			id          uuid.UUID
			deviceID    uuid.UUID
			commandName string
			payload     []byte
			status      string
			executedAt  *time.Time
			createdAt   time.Time
		)

		if err := rows.Scan(&id, &deviceID, &commandName, &payload, &status, &executedAt, &createdAt); err != nil {
			return nil, nil, fmt.Errorf("failed to scan commands: %w", err)
		}

		t, err := command.RehydrateCommand(
			id,
			deviceID,
			commandName,
			payload,
			status,
			executedAt,
			createdAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to rehydrate command: %w", err)
		}

		result = append(result, t)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("rows error: %w", err)
	}

	var nextCur *pagination.Cursor
	if len(result) > limit {
		lastVisible := result[limit-1]
		result = result[:limit]
		nextCur = pagination.NewCursor(uuid.UUID(lastVisible.ID), lastVisible.CreatedAt)
	}

	return result, nextCur, nil
}
