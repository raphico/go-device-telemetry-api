package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/raphico/go-device-telemetry-api/internal/config"
	"github.com/raphico/go-device-telemetry-api/internal/domain/token"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

type UserHandler struct {
	log          *logger.Logger
	cfg          config.Config
	userService  *user.Service
	tokenService *token.Service
}

func NewUserHandler(
	log *logger.Logger,
	cfg config.Config,
	userService *user.Service,
	tokenService *token.Service,
) *UserHandler {
	return &UserHandler{
		log:          log,
		cfg:          cfg,
		userService:  userService,
		tokenService: tokenService,
	}
}

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Email, username, and password are required")
		return
	}

	u, err := h.userService.RegisterUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmail):
			WriteJSONError(w, http.StatusBadRequest, invalidEmail, "Invalid email")
			return
		case errors.Is(err, user.ErrInvalidUsername):
			WriteJSONError(w, http.StatusBadRequest, invalidUsername, "Invalid username")
			return
		case errors.Is(err, user.ErrInvalidPassword):
			WriteJSONError(w, http.StatusBadRequest, invalidPassword, "Invalid password")
			return
		case errors.Is(err, user.ErrEmailAlreadyExists):
			WriteJSONError(w, http.StatusConflict, emailExists, "Email already exists")
			return
		case errors.Is(err, user.ErrUsernameAlreadyExists):
			WriteJSONError(w, http.StatusConflict, usernameExists, "Username already exists")
			return
		default:
			h.log.Error(fmt.Sprintf("failed to register user: %v", err))
			WriteJSONError(w, http.StatusInternalServerError, internalError, "Internal server error")
			return
		}
	}

	res := registerUserResponse{
		ID:       u.ID.String(),
		Username: u.Username.String(),
		Email:    u.Email.String(),
	}

	WriteJSON(w, http.StatusCreated, res, nil)
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Email and password are required")
		return
	}

	user, err := h.userService.AuthenticateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, invalidCredentials, "Invalid email or password")
		return
	}

	accessToken, err := h.tokenService.GenerateAccessToken(user.ID)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, internalError, "Failed to generate access token")
		return
	}

	refreshToken, err := h.tokenService.CreateRefreshToken(r.Context(), user.ID)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, internalError, "Failed to create refresh token")
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Plaintext,
		Path:     "/",
		HttpOnly: true,
		Expires:  refreshToken.ExpiresAt,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.Env == "production",
	}
	http.SetCookie(w, cookie)

	resp := tokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(h.cfg.AccessTokenTTL.Seconds()),
	}

	WriteJSON(w, http.StatusOK, resp, nil)
}

func (h *UserHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, invalidRequest, "Refresh token missing")
		return
	}

	refreshTok := cookie.Value
	accessToken, refreshToken, err := h.tokenService.RotateTokens(r.Context(), refreshTok)
	if err != nil {
		switch {
		case errors.Is(err, token.ErrTokenNotFound):
			WriteJSONError(w, http.StatusUnauthorized, invalidGrant, "Invalid or expired refresh token")
		default:
			h.log.Error(fmt.Sprintf("failed to refresh access token: %v", err))
			WriteJSONError(w, http.StatusInternalServerError, internalError, "Internal server error")
		}
		return
	}

	cookie = &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Plaintext,
		Path:     "/",
		HttpOnly: true,
		Expires:  refreshToken.ExpiresAt,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.Env == "production",
	}
	http.SetCookie(w, cookie)

	resp := tokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(h.cfg.AccessTokenTTL.Seconds()),
	}

	WriteJSON(w, http.StatusOK, resp, nil)
}

func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	refreshToken := cookie.Value
	err = h.tokenService.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, token.ErrTokenNotFound):
			w.WriteHeader(http.StatusNoContent)
		default:
			h.log.Error(fmt.Sprintf("failed to revoke refresh token: %v", err))
			WriteJSONError(w, http.StatusInternalServerError, internalError, "Logout failed, please try again")
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		Secure:   h.cfg.Env == "production",
	})

	w.WriteHeader(http.StatusNoContent)
}
