package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
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

	name, err := device.NewName(req.Name)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	status, err := device.NewStatus(req.Status)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	deviceType, err := device.NewDeviceType(req.DeviceType)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	if req.Metadata == nil {
		req.Metadata = make(map[string]any)
	}

	dev, err := h.device.CreateDevice(r.Context(), userId, name, status, deviceType, req.Metadata)
	if err != nil {
		switch {
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
	userId, ok := GetUserID(r.Context())
	if !ok {
		h.log.Debug(fmt.Sprint("missing user id in context", "path", r.URL.Path))
		WriteUnauthorizedError(w)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := device.NewDeviceID(idStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid device id")
		return
	}

	dev, err := h.device.GetDevice(r.Context(), id, userId)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrDeviceNotFound):
			WriteJSONError(w, http.StatusNotFound, notfound, "device not found")
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

type listDevicesMeta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	Limit      int    `json:"limit"`
}

func (h *DeviceHandler) HandleListDevices(w http.ResponseWriter, r *http.Request) {
	userId, ok := GetUserID(r.Context())
	if !ok {
		h.log.Debug(fmt.Sprint("missing user id in context", "path", r.URL.Path))
		WriteUnauthorizedError(w)
		return
	}

	limit := pagination.DefaultLimit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err != nil {
			WriteJSONError(w, http.StatusBadRequest, invalidRequest, "limit must be a positive integer")
			return
		} else {
			limit = pagination.ClampLimit(v)
		}
	}

	var cur *pagination.Cursor
	if cstr := r.URL.Query().Get("cursor"); cstr != "" {
		if decoded, err := pagination.Decode(cstr); err != nil {
			WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
			return
		} else {
			cur = &decoded
		}
	}

	devs, next, err := h.device.ListUserDevices(r.Context(), userId, limit, cur)
	if err != nil {
		WriteInternalError(w)
		return
	}

	out := make([]deviceResponse, 0, len(devs))
	for _, d := range devs {
		out = append(out, deviceResponse{
			ID:         d.ID.String(),
			Name:       d.Name.String(),
			DeviceType: string(d.DeviceType),
			Status:     string(d.Status),
			Metadata:   d.Metadata,
		})
	}

	var nextStr string
	if next != nil {
		s := pagination.Encode(*next)
		nextStr = s
	}

	meta := listDevicesMeta{
		NextCursor: nextStr,
		Limit:      limit,
	}

	WriteJSON(w, http.StatusOK, out, meta)
}

type updateDeviceRequest struct {
	Name       string         `json:"name"`
	DeviceType string         `json:"device_type"`
	Metadata   map[string]any `json:"metadata"`
}

func (h *DeviceHandler) HandleUpdateDevice(w http.ResponseWriter, r *http.Request) {
	userId, ok := GetUserID(r.Context())
	if !ok {
		h.log.Debug(fmt.Sprint("missing user id in context", "path", r.URL.Path))
		WriteUnauthorizedError(w)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := device.NewDeviceID(idStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid device id")
		return
	}

	var req updateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Invalid request body")
		return
	}

	if req.Name == "" && req.DeviceType == "" && req.Metadata == nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "at least one field must be provided")
		return
	}

	update := device.UpdateDeviceInput{}
	if req.Name != "" {
		n, err := device.NewName(req.Name)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
			return
		}
		update.Name = &n
	}

	if req.DeviceType != "" {
		dt, err := device.NewDeviceType(req.DeviceType)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
			return
		}
		update.DeviceType = &dt
	}

	if req.Metadata != nil {
		update.Metadata = &req.Metadata
	}

	dev, err := h.device.UpdateDevice(r.Context(), id, userId, update)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrDeviceNotFound):
			WriteJSONError(w, http.StatusNotFound, notfound, "device not found")
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
