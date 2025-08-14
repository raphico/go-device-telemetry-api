package db

import (
	"database/sql"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func Migrate(connString string) error {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return err
	}

	defer db.Close()

	goose.SetBaseFS(MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
