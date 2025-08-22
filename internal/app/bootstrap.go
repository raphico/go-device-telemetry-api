package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raphico/go-device-telemetry-api/internal/auth"
	"github.com/raphico/go-device-telemetry-api/internal/command"
	"github.com/raphico/go-device-telemetry-api/internal/config"
	"github.com/raphico/go-device-telemetry-api/internal/db"
	"github.com/raphico/go-device-telemetry-api/internal/device"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
	"github.com/raphico/go-device-telemetry-api/internal/telemetry"
	"github.com/raphico/go-device-telemetry-api/internal/token"
	transporthttp "github.com/raphico/go-device-telemetry-api/internal/transport/http"
	"github.com/raphico/go-device-telemetry-api/internal/user"
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

	commandRepo := db.NewCommandRepository(dbpool)
	commandService := command.NewService(commandRepo)
	commandHandler := transporthttp.NewCommandHandler(log, commandService)

	userMiddleware := transporthttp.NewUserMiddleware(tokenService)

	router := transporthttp.NewRouter(
		log,
		userMiddleware,
		authHandler,
		deviceHandler,
		telemetryHandler,
		commandHandler,
	)

	return router
}
