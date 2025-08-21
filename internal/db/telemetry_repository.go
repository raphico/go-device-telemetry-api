package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
	"github.com/raphico/go-device-telemetry-api/internal/domain/telemetry"
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
		INSERT INTO telemetry (device_id, type, payload)	
		VALUES ($1, $2, $3)
		RETURNING id, recorded_at
	`

	err = r.db.QueryRow(ctx, query, t.DeviceID, t.TelemetryType, jsonPayload).Scan(&t.ID, &t.RecordedAt)

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
