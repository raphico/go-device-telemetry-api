package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

type DeviceHandler struct {
	log     *logger.Logger
	service *device.Service
}

func NewDeviceHandler(log *logger.Logger, deviceService *device.Service) *DeviceHandler {
	return &DeviceHandler{
		log:     log,
		service: deviceService,
	}
}

type createDeviceRequest struct {
	Name       string         `json:"name"`
	DeviceType string         `json:"device_type"`
	Status     string         `json:"status"`
	Metadata   map[string]any `json:"metadata"`
}

type createDeviceResponse struct {
	ID string `json:"id"`
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

	d, err := h.service.CreateDevice(r.Context(), userId, req.Name, req.Status, req.DeviceType, req.Metadata)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrNameRequired),
			errors.Is(err, device.ErrNameTooShort),
			errors.Is(err, device.ErrNameTooLong),
			errors.Is(err, device.ErrNameInvalidChars),
			errors.Is(err, device.ErrInvalidStatus),
			errors.Is(err, device.ErrInvalidDeviceType),
			errors.Is(err, device.ErrInvalidMetadata):
			WriteJSONError(w, http.StatusBadRequest, validationError, err.Error())
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

	res := createDeviceResponse{
		ID: d.ID.String(),
		Name: d.Name.String(),
		DeviceType: string(d.DeviceType),
		Status: string(d.Status),
		Metadata: d.Metadata,
	}

	WriteJSON(w, http.StatusCreated, res, nil)
}
