package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorCode string

const (
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidPayload ErrorCode = "INVALID_PAYLOAD"
)

type ErrorEnvelope struct {
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func Error(w http.ResponseWriter, status int, message string, code ErrorCode) {
	payload := ErrorEnvelope{
		Error: ErrorResponse{
			Message: message,
			Code:    code,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode JSON error response", "error", err)
	}
}

func ValidationError(w http.ResponseWriter, status int, message string, code ErrorCode, fields map[string]string) {
	payload := ErrorEnvelope{
		Error: ErrorResponse{
			Message: message,
			Code:    code,
			Fields:  fields,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode JSON validation error response", "error", err)
	}
}
