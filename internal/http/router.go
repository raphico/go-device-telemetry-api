package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

func NewRouter(
	log *logger.Logger,
	userMw *UserMiddleware,
	authHandler *AuthHandler,
	deviceHandler *DeviceHandler,
	telemetryHandler *TelemetryHandler,
	commandHandler *CommandHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(userMw.AuthMiddleware)
	r.Use(LoggingMiddleware(log))
	r.Use(chimw.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.HandleRegisterUser)
			r.Post("/login", authHandler.HandleLoginUser)
			r.Post("/refresh", authHandler.HandleRefreshAccessToken)
		})

		r.Group(func(r chi.Router) {
			r.Use(userMw.RequireAuthMiddleware)

			r.Post("/auth/logout", authHandler.HandleLogoutUser)

			r.Route("/devices", func(r chi.Router) {
				r.Post("/", deviceHandler.HandleCreateDevice)
				r.Get("/", deviceHandler.HandleListDevices)
				r.Get("/{device_id}", deviceHandler.HandleGetDevice)
				r.Post("/{device_id}", deviceHandler.HandleUpdateDevice)

				r.Route("/{device_id}/telemetry", func(r chi.Router) {
					r.Post("/", telemetryHandler.HandleCreateTelemetry)
					r.Get("/", telemetryHandler.HandleGetDeviceTelemetry)
				})

				r.Route("/{device_id}/commands", func(r chi.Router) {
					r.Post("/", commandHandler.HandleCreateCommand)
				})
			})
		})

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, "OK", nil)
		})
	})

	return r
}
