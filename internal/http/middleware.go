package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/raphico/go-device-telemetry-api/internal/domain/token"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

type UserMiddleware struct {
	tokenService *token.Service
}

type contextKey struct{}

var userCtxKey = &contextKey{}

func NewUserMiddleware(tokenService *token.Service) *UserMiddleware {
	return &UserMiddleware{
		tokenService: tokenService,
	}
}

func (um *UserMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		userId, err := um.tokenService.ValidateAccessToken(accessToken)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (um *UserMiddleware) RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(userCtxKey).(user.UserID)
		if !ok {
			WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "authentication required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserID(ctx context.Context) (user.UserID, bool) {
	userID, ok := ctx.Value(userCtxKey).(user.UserID)
	return userID, ok
}

func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
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
