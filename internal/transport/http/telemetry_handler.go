package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/device"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
	"github.com/raphico/go-device-telemetry-api/internal/telemetry"
)

type TelemetryHandler struct {
	log       *logger.Logger
	telemetry *telemetry.Service
}

func NewTelemetryHandler(log *logger.Logger, telemetry *telemetry.Service) *TelemetryHandler {
	return &TelemetryHandler{
		log:       log,
		telemetry: telemetry,
	}
}

type createTelemetryRequest struct {
	TelemetryType string `json:"telemetry_type"`
	Payload       any    `json:"payload"`
	RecordedAt    string `json:"recorded_at"`
}

type telemetryResponse struct {
	ID            string         `json:"id"`
	TelemetryType string         `json:"telemetry_type"`
	Payload       map[string]any `json:"payload"`
	RecordedAt    time.Time      `json:"recorded_at"`
}

func (h *TelemetryHandler) HandleCreateTelemetry(w http.ResponseWriter, r *http.Request) {
	deviceIDStr := chi.URLParam(r, "device_id")
	deviceID, err := device.NewDeviceID(deviceIDStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid device id")
		return
	}

	var req createTelemetryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid request body")
		return
	}

	telemetryType, err := telemetry.NewTelemetryType(req.TelemetryType)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	payload, err := telemetry.NewPayload(req.Payload)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	recordedAt, err := telemetry.NewRecordedAt(req.RecordedAt)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	t, err := h.telemetry.CreateTelemetry(r.Context(), deviceID, telemetryType, payload, recordedAt)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrDeviceNotFound):
			WriteJSONError(w, http.StatusNotFound, notfound, "device not found")
		default:
			h.log.Error(fmt.Sprintf("failed to add telemetry: %v", err))
			WriteInternalError(w)
		}
		return
	}

	res := telemetryResponse{
		ID:            t.ID.String(),
		TelemetryType: t.TelemetryType.String(),
		Payload:       t.Payload,
		RecordedAt:    t.RecordedAt.Time(),
	}

	WriteJSON(w, http.StatusCreated, res, nil)
}

func (h *TelemetryHandler) HandleGetDeviceTelemetry(w http.ResponseWriter, r *http.Request) {
	deviceIDStr := chi.URLParam(r, "device_id")
	deviceID, err := device.NewDeviceID(deviceIDStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid device id")
		return
	}

	limit := pagination.DefaultLimit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err != nil || v < 0 {
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

	telemetry, next, err := h.telemetry.ListDeviceTelemetry(r.Context(), deviceID, limit, cur)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to get device telemetry: %v", err))
		WriteInternalError(w)
		return
	}

	out := make([]telemetryResponse, 0, len(telemetry))
	for _, t := range telemetry {
		out = append(out, telemetryResponse{
			ID:            t.ID.String(),
			TelemetryType: t.TelemetryType.String(),
			Payload:       t.Payload,
			RecordedAt:    t.RecordedAt.Time(),
		})
	}

	var nextStr string
	if next != nil {
		nextStr = pagination.Encode(*next)
	}

	meta := pageMeta{
		NextCursor: nextStr,
		Limit:      limit,
	}

	WriteJSON(w, http.StatusOK, out, meta)
}
