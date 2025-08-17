package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/config"
	"github.com/raphico/go-device-telemetry-api/internal/db"
	"github.com/raphico/go-device-telemetry-api/internal/domain/token"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
	transporthttp "github.com/raphico/go-device-telemetry-api/internal/transport/http"
)

func BuildApp(log *logger.Logger, dbpool *pgxpool.Pool, cfg config.Config) http.Handler {
	tokenRepo := db.NewTokenRepository(dbpool)
	jwtGenerator := token.NewJWTAdapter([]byte(cfg.JWTSecret))
	tokenService := token.NewService(jwtGenerator, tokenRepo, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	userRepo := db.NewUserRepository(dbpool)
	userService := user.NewService(userRepo)
	userHandler := transporthttp.NewUserHandler(log, cfg, userService, tokenService)

	router := transporthttp.NewRouter(log, userHandler)

	return router
}
