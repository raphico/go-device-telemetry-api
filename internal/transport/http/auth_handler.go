package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/raphico/go-device-telemetry-api/internal/auth"
	"github.com/raphico/go-device-telemetry-api/internal/config"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
	"github.com/raphico/go-device-telemetry-api/internal/token"
	"github.com/raphico/go-device-telemetry-api/internal/user"
)

type AuthHandler struct {
	log  *logger.Logger
	cfg  config.Config
	auth *auth.Service
}

func NewAuthHandler(
	log *logger.Logger,
	cfg config.Config,
	authService *auth.Service,
) *AuthHandler {
	return &AuthHandler{
		log:  log,
		cfg:  cfg,
		auth: authService,
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

func (h *AuthHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Invalid request body")
		return
	}

	username, err := user.NewUsername(req.Username)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	email, err := user.NewEmail(req.Email)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	password, err := user.NewPassword(req.Password)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	u, err := h.auth.Register(r.Context(), username, email, password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailAlreadyExists),
			errors.Is(err, user.ErrUsernameTaken):
			WriteJSONError(w, http.StatusConflict, conflict, err.Error())
		default:
			h.log.Error(fmt.Sprintf("failed to register user: %v", err))
			WriteInternalError(w)
		}
		return
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

func (h *AuthHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Invalid request body")
		return
	}

	email, err := user.NewEmail(req.Email)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.auth.Login(r.Context(), email, req.Password)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			WriteJSONError(w, http.StatusUnauthorized, invalidRequest, "Invalid credentials")
			return
		}

		h.log.Error(fmt.Sprintf("failed to authenticate user: %v", err))
		WriteInternalError(w)
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

func (h *AuthHandler) HandleRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Refresh token missing")
		return
	}

	refreshTok := cookie.Value
	accessToken, refreshToken, err := h.auth.Refresh(r.Context(), refreshTok)
	if err != nil {
		switch {
		case errors.Is(err, token.ErrTokenNotFound):
			WriteJSONError(w, http.StatusUnauthorized, unauthorized, "Invalid or expired refresh token")
		default:
			h.log.Error(fmt.Sprintf("failed to refresh access token: %v", err))
			WriteInternalError(w)
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

func (h *AuthHandler) HandleLogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	refreshToken := cookie.Value
	err = h.auth.Logout(r.Context(), refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, token.ErrTokenNotFound):
			w.WriteHeader(http.StatusNoContent)
		default:
			h.log.Error(fmt.Sprintf("failed to revoke refresh token: %v", err))
			WriteInternalError(w)
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
