package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

func NewRouter(log *logger.Logger, userMw *UserMiddleware, userHandler *UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(userMw.AuthMiddleware)
	r.Use(loggingMiddleware(log))
	r.Use(chimw.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.RegisterUser)
			r.Post("/login", userHandler.LoginUser)
			r.Post("/refresh", userHandler.RefreshAccessToken)
		})

		r.Group(func(r chi.Router) {
			r.Use(userMw.RequireAuthMiddleware)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			})
		})

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				log.Error(fmt.Sprintf("failed to write health response: %v", err))
			}
		})

	})

	return r
}
