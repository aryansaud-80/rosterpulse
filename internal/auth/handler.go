package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/aryansaud-80/rosterpulse/internal/shared/httpx"
	"github.com/aryansaud-80/rosterpulse/internal/shared/validator"
)

type AuthHandler struct {
	service IAuthService
	log     *slog.Logger
	env     string
}

func NewAuthHandler(service IAuthService, log *slog.Logger, env string) *AuthHandler {
	return &AuthHandler{
		service: service,
		env:     env,
		log:     log.With("component", "auth.handler"),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode login request", "error", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid request payload", httpx.ErrCodeInvalidPayload)
		return
	}

	if err := validator.Validate.Struct(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Validation failed", httpx.ErrCodeValidation)
		return
	}

	aToken, rToken, uRep, err := h.service.Login(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrUserNotFound):
			httpx.Error(w, http.StatusUnauthorized, "invalid email or password", httpx.ErrCodeUnauthorized)

		case errors.Is(err, ErrTokenGenerationFailed), errors.Is(err, ErrSessionCreation):
			httpx.Error(w, http.StatusInternalServerError, "internal server error", httpx.ErrCodeInternal)

		default:
			httpx.Error(w, http.StatusInternalServerError, "internal server error", httpx.ErrCodeInternal)
		}
		return
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    rToken,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   h.env == "production",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	http.SetCookie(w, &cookie)
	response := &LoginResponse{
		AccessToken: aToken,
		User:        uRep,
	}

	httpx.JSON(w, http.StatusOK, response)
}
