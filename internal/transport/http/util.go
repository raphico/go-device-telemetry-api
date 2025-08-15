package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raphico/go-device-telemetry-api/internal/logger"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteJSONError(w http.ResponseWriter, status int, code, msg string, log *logger.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorResponse{
		Error:   code,
		Message: msg,
	}); err != nil {
		log.Error(fmt.Sprintf("failed to write JSON error: %v", err))
	}
}
