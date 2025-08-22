package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/command"
	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/device"
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

	if err := c.Status.SetStatus(status); err != nil {
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

func (r *CommandRepository) FindById(
	ctx context.Context,
	id command.CommandID,
	deviceID device.DeviceID,
) (*command.Command, error) {
	var (
		commandID   uuid.UUID
		dbDeviceID  uuid.UUID
		commandName string
		payload     []byte
		status      string
		executedAt  *time.Time
		createdAt   time.Time
	)

	query := `
    	SELECT id, device_id, command_name, payload, status, executed_at, created_at
		FROM commands
		WHERE id = $1 AND device_id = $2
	`

	err := r.db.QueryRow(ctx, query, id, deviceID).Scan(
		&commandID,
		&dbDeviceID,
		&commandName,
		&payload,
		&status,
		&executedAt,
		&createdAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, command.ErrCommandNotFound
		}

		return nil, fmt.Errorf("failed to find command by id: %w", err)
	}

	return command.RehydrateCommand(
		commandID,
		dbDeviceID,
		commandName,
		payload,
		status,
		executedAt,
		createdAt,
	)
}

func (r *CommandRepository) UpdateStatus(ctx context.Context, c *command.Command) error {
	if !c.ExecutedAt.Valid() {
		return fmt.Errorf("invalid executed_at")
	}

	query := `
		UPDATE commands
		SET status = $1, executed_at = $2
		WHERE id = $3 AND device_id = $4
	`

	tag, err := r.db.Exec(
		ctx,
		query,
		c.Status.String(),
		c.ExecutedAt.Time(),
		c.ID,
		c.DeviceID,
	)

	if err != nil {
		return fmt.Errorf("failed to update command: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return command.ErrCommandNotFound
	}

	return nil
}
