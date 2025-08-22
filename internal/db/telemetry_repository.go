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
	"github.com/raphico/go-device-telemetry-api/internal/device"
	"github.com/raphico/go-device-telemetry-api/internal/telemetry"
)

type TelemetryRepository struct {
	db *pgxpool.Pool
}

func NewTelemetryRepository(db *pgxpool.Pool) *TelemetryRepository {
	return &TelemetryRepository{
		db: db,
	}
}

func (r *TelemetryRepository) Create(ctx context.Context, t *telemetry.Telemetry) error {
	jsonPayload, err := json.Marshal(t.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO telemetry (device_id, telemetry_type, payload)	
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err = r.db.QueryRow(ctx, query, t.DeviceID, t.TelemetryType, jsonPayload).Scan(&t.ID, &t.CreatedAt)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == "23503" && pgError.ConstraintName == "telemetry_device_id_fkey" {
				return device.ErrDeviceNotFound
			}
		}
		return fmt.Errorf("failed to insert telemetry: %w", err)
	}

	return nil
}

func (r *TelemetryRepository) FindTelemetry(
	ctx context.Context,
	deviceID device.DeviceID,
	limit int,
	cursor *pagination.Cursor,
) ([]*telemetry.Telemetry, *pagination.Cursor, error) {
	var query string
	var args []any

	if cursor == nil {
		query = `
    		SELECT id, device_id, telemetry_type, payload, recorded_at, created_at
    		FROM telemetry
    		WHERE device_id = $1
    		ORDER BY created_at ASC, id ASC
    		LIMIT $2
    	`
		args = []any{deviceID, limit + 1}
	} else {
		query = `
    		SELECT id, device_id, telemetry_type, payload, recorded_at, created_at
    		FROM telemetry
    		WHERE device_id = $1
				AND (created_at, id) > ($2, $3)
    		ORDER BY created_at ASC, id ASC
    		LIMIT $4
		`
		args = []any{deviceID, cursor.CreatedAt, cursor.ID, limit + 1}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query telemetry: %w", err)
	}

	var result []*telemetry.Telemetry
	for rows.Next() {
		var (
			id            uuid.UUID
			deviceID      uuid.UUID
			telemetryType string
			payload       []byte
			recordedAt    time.Time
			createdAt     time.Time
		)

		if err := rows.Scan(&id, &deviceID, &telemetryType, &payload, &recordedAt, &createdAt); err != nil {
			return nil, nil, fmt.Errorf("failed to scan telemetry: %w", err)
		}

		t, err := telemetry.RehydrateTelemetry(
			id,
			deviceID,
			telemetryType,
			payload,
			recordedAt,
			createdAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to rehydrate telemetry: %w", err)
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
