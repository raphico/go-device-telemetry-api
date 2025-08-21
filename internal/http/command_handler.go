package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/raphico/go-device-telemetry-api/internal/domain/command"
	"github.com/raphico/go-device-telemetry-api/internal/domain/device"
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

type createCommandResponse struct {
	ID          string `json:"id"`
	CommandName string `json:"command_name"`
	Payload     any    `json:"payload"`
	Status      string `json:"status"`
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

	c, err := h.command.CreateCommand(r.Context(), deviceID, commandName, payload)
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

	res := createCommandResponse{
		ID:          c.ID.String(),
		CommandName: c.Name.String(),
		Payload:     c.Payload,
		Status:      c.Status.String(),
	}

	WriteJSON(w, http.StatusCreated, res, nil)
}
