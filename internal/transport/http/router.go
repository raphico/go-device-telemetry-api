package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

func NewRouter(log *logger.Logger, userHandler *UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(loggingMiddleware(log))
	r.Use(chimw.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("OK")); err != nil {
				log.Error(fmt.Sprintf("failed to write health response: %v", err))
			}
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", userHandler.RegisterUser)
		})
	})

	return r
}

func loggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info(fmt.Sprintf(
				"HTTP request: method=%s, path=%s, remote=%s, duration=%s, reqID=%s",
				r.Method, r.URL.Path, r.RemoteAddr, time.Since(start), chimw.GetReqID(r.Context()),
			))
		})
	}
}
