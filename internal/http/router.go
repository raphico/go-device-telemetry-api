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
			r.Post("/register", authHandler.RegisterUser)
			r.Post("/login", authHandler.LoginUser)
			r.Post("/refresh", authHandler.RefreshAccessToken)
		})

		r.Group(func(r chi.Router) {
			r.Use(userMw.RequireAuthMiddleware)

			r.Post("/auth/logout", authHandler.LogoutUser)

			r.Route("/devices", func(r chi.Router) {
				r.Post("/", deviceHandler.HandleCreateDevice)
				r.Get("/{id}", deviceHandler.HandleGetDevice)
			})
		})

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, "OK", nil)
		})
	})

	return r
}
