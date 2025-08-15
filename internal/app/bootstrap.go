package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/db"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
	transporthttp "github.com/raphico/go-device-telemetry-api/internal/transport/http"
)

func BuildApp(log *logger.Logger, dbpool *pgxpool.Pool) http.Handler {
	userRepo := db.NewUserRepository(dbpool)
	userService := user.NewService(userRepo)
	userHandler := transporthttp.NewUserHandler(log, userService)

	router := transporthttp.NewRouter(log, userHandler)

	return router
}
