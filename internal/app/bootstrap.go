package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/config"
	"github.com/raphico/go-device-telemetry-api/internal/db"
	"github.com/raphico/go-device-telemetry-api/internal/domain/auth"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
	"github.com/raphico/go-device-telemetry-api/internal/domain/telemetry"
	"github.com/raphico/go-device-telemetry-api/internal/domain/token"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
	transporthttp "github.com/raphico/go-device-telemetry-api/internal/http"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

func BuildApp(log *logger.Logger, dbpool *pgxpool.Pool, cfg config.Config) http.Handler {
	tokenRepo := db.NewTokenRepository(dbpool)
	jwtGenerator := token.NewJWTAdapter([]byte(cfg.JWTSecret))
	tokenService := token.NewService(jwtGenerator, tokenRepo, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	userRepo := db.NewUserRepository(dbpool)
	userService := user.NewService(userRepo)

	authService := auth.NewService(userService, tokenService)
	authHandler := transporthttp.NewAuthHandler(log, cfg, authService)

	deviceRepo := db.NewDeviceRepository(dbpool)
	deviceService := device.NewService(deviceRepo)
	deviceHandler := transporthttp.NewDeviceHandler(log, deviceService)

	telemetryRepo := db.NewTelemetryRepository(dbpool)
	telemetryService := telemetry.NewService(telemetryRepo)
	telemetryHandler := transporthttp.NewTelemetryHandler(log, telemetryService)

	userMiddleware := transporthttp.NewUserMiddleware(tokenService)

	router := transporthttp.NewRouter(log, userMiddleware, authHandler, deviceHandler, telemetryHandler)

	return router
}
