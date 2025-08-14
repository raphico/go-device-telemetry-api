package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func Migrate(connString string, log *logger.Logger) error {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error(fmt.Sprintf("failed to close database connection, %v", err))
		}
	}()

	goose.SetBaseFS(MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
