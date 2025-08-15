package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

type UserHandler struct {
	log     *logger.Logger
	service *user.Service
}

func NewUserHandler(log *logger.Logger, service *user.Service) *UserHandler {
	return &UserHandler{
		log:     log,
		service: service,
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	u, err := h.service.RegisterUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmail):
			WriteJSONError(w, http.StatusBadRequest, "invalid_email", err.Error(), h.log)
		case errors.Is(err, user.ErrInvalidUsername):
			WriteJSONError(w, http.StatusBadRequest, "invalid_username", err.Error(), h.log)
		case errors.Is(err, user.ErrInvalidPassword):
			WriteJSONError(w, http.StatusBadRequest, "invalid_password", err.Error(), h.log)
		case errors.Is(err, user.ErrEmailAlreadyExists):
			WriteJSONError(w, http.StatusConflict, "email_exists", err.Error(), h.log)
		case errors.Is(err, user.ErrUsernameAlreadyExits):
			WriteJSONError(w, http.StatusConflict, "username_exists", err.Error(), h.log)
		default:
			h.log.Error(fmt.Sprintf("failed to register user: %v", err))
			WriteJSONError(w, http.StatusInternalServerError, "internal_error", "internal server error", h.log)
		}
		return
	}

	resp := registerUserResponse{
		ID:       string(u.ID),
		Username: u.Username.String(),
		Email:    u.Email.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(fmt.Sprintf("failed to encode response: %v", err))
	}
}
