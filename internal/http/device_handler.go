package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

type DeviceHandler struct {
	log    *logger.Logger
	device *device.Service
}

func NewDeviceHandler(log *logger.Logger, deviceService *device.Service) *DeviceHandler {
	return &DeviceHandler{
		log:    log,
		device: deviceService,
	}
}

type createDeviceRequest struct {
	Name       string         `json:"name"`
	DeviceType string         `json:"device_type"`
	Status     string         `json:"status"`
	Metadata   map[string]any `json:"metadata"`
}

type deviceResponse struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	DeviceType string         `json:"device_type"`
	Status     string         `json:"status"`
	Metadata   map[string]any `json:"metadata"`
}

func (h *DeviceHandler) HandleCreateDevice(w http.ResponseWriter, r *http.Request) {
	var req createDeviceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Invalid request body")
		return
	}

	userId, ok := GetUserID(r.Context())
	if !ok {
		h.log.Debug(fmt.Sprint("missing user id in context", "path", r.URL.Path))
		WriteUnauthorizedError(w)
	}

	dev, err := h.device.CreateDevice(r.Context(), userId, req.Name, req.Status, req.DeviceType, req.Metadata)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrNameRequired),
			errors.Is(err, device.ErrNameTooShort),
			errors.Is(err, device.ErrNameTooLong),
			errors.Is(err, device.ErrNameInvalidChars),
			errors.Is(err, device.ErrInvalidStatus),
			errors.Is(err, device.ErrInvalidDeviceType),
			errors.Is(err, device.ErrInvalidMetadata):
			WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
			return

		case errors.Is(err, user.ErrUserNotFound):
			WriteJSONError(w, http.StatusUnauthorized, unauthorized, "User does not exist")
			return

		default:
			h.log.Error(fmt.Sprint("failed to create device", "error", err))
			WriteInternalError(w)
			return
		}
	}

	res := deviceResponse{
		ID:         dev.ID.String(),
		Name:       dev.Name.String(),
		DeviceType: string(dev.DeviceType),
		Status:     string(dev.Status),
		Metadata:   dev.Metadata,
	}

	WriteJSON(w, http.StatusCreated, res, nil)
}

func (h *DeviceHandler) HandleGetDevice(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid device id")
		return
	}

	userId, ok := GetUserID(r.Context())
	if !ok {
		h.log.Debug(fmt.Sprint("missing user id in context", "path", r.URL.Path))
		WriteUnauthorizedError(w)
		return
	}

	dev, err := h.device.GetDevice(r.Context(), device.DeviceID(id), userId)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrDeviceNotFound):
			WriteJSONError(w, http.StatusBadRequest, invalidRequest, "device not found")
		default:
			h.log.Error(fmt.Sprintf("failed to get device: %v", err))
			WriteInternalError(w)
		}
		return
	}

	res := deviceResponse{
		ID:         dev.ID.String(),
		Name:       dev.Name.String(),
		DeviceType: string(dev.DeviceType),
		Status:     string(dev.Status),
		Metadata:   dev.Metadata,
	}

	WriteJSON(w, http.StatusOK, res, nil)
}
