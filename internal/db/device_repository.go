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
	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
)

type DeviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) Create(ctx context.Context, dev *device.Device) error {
	jsonMetadata, err := json.Marshal(dev.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO devices (user_id, name, device_type, status, metadata)
		VALUES ($1, $2, $3, $4, $5)
		returning id, created_at, updated_at
	`

	err = r.db.QueryRow(
		ctx,
		query,
		dev.UserID,
		dev.Name,
		dev.DeviceType,
		dev.Status,
		jsonMetadata,
	).Scan(&dev.ID, &dev.CreatedAt, &dev.UpdatedAt)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == "23503" && pgError.ConstraintName == "devices_user_id_fkey" {
				return user.ErrUserNotFound
			}
		}

		return fmt.Errorf("failed to insert device: %w", err)
	}

	return nil
}

func (r *DeviceRepository) FindById(
	ctx context.Context,
	id device.DeviceID,
	userID user.UserID,
) (*device.Device, error) {
	var (
		deviceID   uuid.UUID
		dbUserID     uuid.UUID
		name       string
		deviceType string
		status     string
		metadata   []byte
		createdAt  time.Time
		updatedAt  time.Time
	)

	query := `
		SELECT id, user_id, name, device_type, status, metadata, created_at, updated_at
		FROM devices
		WHERE id = $1 AND user_id = $2
	`

	err := r.db.QueryRow(ctx, query, id, dbUserID).Scan(
		&deviceID,
		&userID,
		&name,
		&deviceType,
		&status,
		&metadata,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, device.ErrDeviceNotFound
		}

		return nil, fmt.Errorf("failed to find device by id: %w", err)
	}

	return device.RehydrateDevice(
		device.DeviceID(deviceID),
		user.UserID(userID),
		name,
		deviceType,
		status,
		metadata,
		createdAt,
		updatedAt,
	)

}

func (r *DeviceRepository) FindDevices(
	ctx context.Context,
	userID user.UserID,
	limit int,
	cursor *pagination.Cursor,
) ([]*device.Device, *pagination.Cursor, error) {
	var (
		query string
		args  []any
	)

	if cursor == nil {
		query = `
			SELECT id, user_id, name, device_type, status, metadata, created_at, updated_at
			FROM devices
			WHERE user_id = $1
			ORDER BY created_at ASC, id ASC
			LIMIT $2
		`
		args = []any{userID, limit + 1}
	} else {
		query = `
			SELECT id, user_id, name, device_type, status, metadata, created_at, updated_at
			FROM devices
			WHERE user_id = $1
			  AND (created_at, id) > ($2, $3)
			ORDER BY created_at ASC, id ASC
			LIMIT $4
		`
		args = []any{userID, cursor.CreatedAt, cursor.ID, limit + 1}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var result []*device.Device

	for rows.Next() {
		var (
			deviceID   uuid.UUID
			uID        uuid.UUID
			name       string
			deviceType string
			status     string
			metadata   []byte
			createdAt  time.Time
			updatedAt  time.Time
		)

		if err := rows.Scan(
			&deviceID,
			&uID,
			&name,
			&deviceType,
			&status,
			&metadata,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, nil, fmt.Errorf("failed to scan device: %w", err)
		}

		dev, err := device.RehydrateDevice(
			device.DeviceID(deviceID),
			user.UserID(uID),
			name,
			deviceType,
			status,
			metadata,
			createdAt,
			updatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to rehydrate device: %w", err)
		}

		result = append(result, dev)
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

func (r *DeviceRepository) Update(ctx context.Context, dev *device.Device) error {
	query := `
		UPDATE devices
        SET name = $1, device_type = $2, metadata = $3, updated_at = NOW()
        WHERE id = $4 AND user_id = $5	
	`

	tag, err := r.db.Exec(ctx, query, dev.Name, dev.DeviceType, dev.Metadata, dev.ID, dev.UserID)
	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return device.ErrDeviceNotFound
	}

	return nil
}
