package http

import (
	"encoding/json"
	"net/http"
)

type errorPayload struct {
	Code    errorCode `json:"code"`
	Message string    `json:"message"`
}

type responseEnvelope struct {
	Success bool          `json:"success"`
	Data    any           `json:"data,omitempty"`
	Meta    any           `json:"meta,omitempty"`
	Error   *errorPayload `json:"error,omitempty"`
}

func WriteJSONError(w http.ResponseWriter, status int, code errorCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := responseEnvelope{
		Success: false,
		Error: &errorPayload{
			Code:    code,
			Message: msg,
		},
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any, meta any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := responseEnvelope{
		Success: true,
		Data:    data,
		Meta:    meta,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
