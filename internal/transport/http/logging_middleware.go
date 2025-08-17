package http

import (
	"fmt"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

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
