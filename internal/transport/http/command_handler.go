package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/raphico/go-device-telemetry-api/internal/command"
	"github.com/raphico/go-device-telemetry-api/internal/common/pagination"
	"github.com/raphico/go-device-telemetry-api/internal/device"
	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

type CommandHandler struct {
	log     *logger.Logger
	command *command.Service
}

func NewCommandHandler(log *logger.Logger, command *command.Service) *CommandHandler {
	return &CommandHandler{
		log:     log,
		command: command,
	}
}

type createCommandRequest struct {
	CommandName string `json:"command_name"`
	Payload     any    `json:"payload"`
}

type commandResponse struct {
	ID          string    `json:"id"`
	CommandName string    `json:"command_name"`
	Payload     any       `json:"payload"`
	Status      string    `json:"status"`
	ExecutedAt  time.Time `json:"executed_at,omitzero"`
}

func (h *CommandHandler) HandleCreateCommand(w http.ResponseWriter, r *http.Request) {
	deviceIDStr := chi.URLParam(r, "device_id")
	deviceID, err := device.NewDeviceID(deviceIDStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid device id")
		return
	}

	var req createCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid request body")
		return
	}

	commandName, err := command.NewName(req.CommandName)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	payload, err := command.NewPayload(req.Payload)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	cmd, err := h.command.CreateCommand(r.Context(), deviceID, commandName, payload)
	if err != nil {
		switch {
		case errors.Is(err, device.ErrDeviceNotFound):
			WriteJSONError(w, http.StatusNotFound, notfound, "device not found")
		default:
			h.log.Error(fmt.Sprintf("failed to add command: %v", err))
			WriteInternalError(w)
		}
		return
	}

	res := commandResponse{
		ID:          cmd.ID.String(),
		CommandName: cmd.Name.String(),
		Payload:     cmd.Payload,
		Status:      cmd.Status.String(),
	}

	if cmd.ExecutedAt.Valid() {
		res.ExecutedAt = cmd.ExecutedAt.Time()
	}

	WriteJSON(w, http.StatusCreated, res, nil)
}

func (h *CommandHandler) HandleGetDeviceCommands(w http.ResponseWriter, r *http.Request) {
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

	commands, next, err := h.command.ListDeviceCommands(r.Context(), deviceID, limit, cur)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to get device commands: %v", err))
		WriteInternalError(w)
		return
	}

	out := make([]commandResponse, 0, len(commands))
	for _, c := range commands {
		cmd := commandResponse{
			ID:          c.ID.String(),
			CommandName: c.Name.String(),
			Payload:     c.Payload,
			Status:      c.Status.String(),
		}
		if c.ExecutedAt.Valid() {
			cmd.ExecutedAt = c.ExecutedAt.Time()
		}
		out = append(out, cmd)
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

type updateCommandStatusRequest struct {
	Status     string `json:"status"`
	ExecutedAt string `json:"executed_at"`
}

func (h *CommandHandler) HandleUpdateCommandStatus(w http.ResponseWriter, r *http.Request) {
	deviceIDStr := chi.URLParam(r, "device_id")
	deviceID, err := device.NewDeviceID(deviceIDStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid device id")
		return
	}

	commandIDStr := chi.URLParam(r, "command_id")
	commandID, err := command.NewCommandID(commandIDStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "invalid command id")
		return
	}

	var req updateCommandStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, "Invalid request body")
		return
	}

	status, err := command.NewStatus(req.Status)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	executedAt, err := command.NewExecutedAt(req.ExecutedAt)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, invalidRequest, err.Error())
		return
	}

	cmd, err := h.command.UpdateCommandStatus(r.Context(), commandID, deviceID, status, executedAt)
	if err != nil {
		switch {
		case errors.Is(err, command.ErrCommandNotFound):
			WriteJSONError(w, http.StatusNotFound, notfound, "command not found")
		default:
			h.log.Error(fmt.Sprintf("failed to update command: %v", err))
			WriteInternalError(w)
		}
		return
	}

	res := commandResponse{
		ID:          cmd.ID.String(),
		CommandName: cmd.Name.String(),
		Payload:     cmd.Payload,
		Status:      cmd.Status.String(),
	}

	if cmd.ExecutedAt.Valid() {
		res.ExecutedAt = cmd.ExecutedAt.Time()
	}

	WriteJSON(w, http.StatusCreated, res, nil)
}
