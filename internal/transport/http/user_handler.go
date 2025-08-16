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
	defer r.Body.Close()
	
	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	u, err := h.service.RegisterUser(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		errorMap := []struct{
			Err error
			HTTP int
			Code string
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
				WriteJSONError(w, e.HTTP, e.Code, e.Message, h.log)
				return;
			}
		}

		h.log.Error(fmt.Sprintf("failed to register user: %v", err))
		WriteJSONError(w, http.StatusInternalServerError, "internal_error", "internal server error", h.log)
		return
	}

	resp := registerUserResponse{
		ID:       u.ID.String(),
		Username: u.Username.String(),
		Email:    u.Email.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(fmt.Sprintf("failed to encode response: %v", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
