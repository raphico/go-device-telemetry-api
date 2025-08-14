package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

func NewRouter(logger *logger.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(loggingMiddleware(logger))
	r.Use(chimw.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
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
