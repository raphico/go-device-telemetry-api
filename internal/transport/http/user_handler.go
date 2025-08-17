package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
		WriteJSONError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		WriteJSONError(w, http.StatusBadRequest, "invalid_request", "email, username, and password are required")
		return
	}

	u, err := h.userService.RegisterUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		errorMap := []struct {
			Err     error
			HTTP    int
			Code    string
			Message string
		}{
			{user.ErrInvalidEmail, http.StatusBadRequest, "invalid_email", "invalid email"},
			{user.ErrInvalidUsername, http.StatusBadRequest, "invalid_username", "invalid username"},
			{user.ErrInvalidPassword, http.StatusBadRequest, "invalid_password", "invalid password"},
			{user.ErrEmailAlreadyExists, http.StatusConflict, "email_exists", "email already exists"},
			{user.ErrUsernameAlreadyExists, http.StatusConflict, "username_exists", "username already exists"},
		}

		for _, e := range errorMap {
			if errors.Is(err, e.Err) {
				WriteJSONError(w, e.HTTP, e.Code, e.Message)
				return
			}
		}

		h.log.Error(fmt.Sprintf("failed to register user: %v", err))
		WriteJSONError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}

	resp := registerUserResponse{
		ID:       u.ID.String(),
		Username: u.Username.String(),
		Email:    u.Email.String(),
	}

	WriteJSON(w, http.StatusCreated, resp)
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
		WriteJSONError(w, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		WriteJSONError(w, http.StatusBadRequest, "invalid_request", "email and password are required")
		return
	}

	user, err := h.userService.AuthenticateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "invalid_credentials", "invalid email or password")
		return
	}

	accessToken, err := h.tokenService.GenerateAccessToken(user.ID)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "internal_error", "failed to generate access token")
		return
	}

	refreshToken, err := h.tokenService.CreateRefreshToken(r.Context(), user.ID)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "internal_error", "failed to create refresh token")
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Plaintext,
		Path:     "/",
		HttpOnly: true,
		Expires:  refreshToken.ExpiresAt,
		SameSite: http.SameSiteLaxMode,
	}
	if h.cfg.Env == "production" {
		cookie.Secure = true
	} else {
		cookie.Secure = false
	}
	http.SetCookie(w, cookie)

	resp := tokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(h.cfg.AccessTokenTTL.Seconds()),
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "invalid_request", "refresh token missing")
		return
	}

	refreshTok := cookie.Value

	accessToken, refreshToken, err := h.tokenService.RotateTokens(r.Context(), refreshTok)
	if err != nil {
		switch {
		case errors.Is(err, token.ErrTokenNotFound):
			WriteJSONError(w, http.StatusUnauthorized, "invalid_grant", "invalid or expired refresh token")
		default:
			h.log.Error(fmt.Sprintf("failed to refresh access token: %v", err))
			WriteJSONError(w, http.StatusInternalServerError, "internal_error", "internal server error")
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
	}
	if h.cfg.Env == "production" {
		cookie.Secure = true
	} else {
		cookie.Secure = false
	}
	http.SetCookie(w, cookie)

	resp := tokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(h.cfg.AccessTokenTTL.Seconds()),
	}

	WriteJSON(w, http.StatusOK, resp)
}
