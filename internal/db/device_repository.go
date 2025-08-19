package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *DeviceRepository) Create(ctx context.Context, device *device.Device) error {
	jsonMetadata, err := json.Marshal(device.Metadata)
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
		device.UserID,
		device.Name,
		device.DeviceType,
		device.Status,
		jsonMetadata,
	).Scan(&device.ID, &device.CreatedAt, &device.UpdatedAt)

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
