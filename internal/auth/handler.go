package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/aryansaud-80/rosterpulse/internal/middleware"
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
		Path:     "/api/auth",
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

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	rawClaims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized", httpx.ErrCodeUnauthorized)
		return
	}

	claims, ok := rawClaims.(*Claims)
	if !ok {
		h.log.Warn("failed to type assert claims from context", "actual_type", "rawClaims", rawClaims)
		httpx.Error(w, http.StatusUnauthorized, "invalid token claims", httpx.ErrCodeUnauthorized)
		return
	}

	userID := claims.UserID
	if strings.TrimSpace(userID) == "" {
		h.log.Warn("user id is empty within valid token claims")
		httpx.Error(w, http.StatusUnauthorized, "invalid token claims: missing user id", httpx.ErrCodeUnauthorized)
		return
	}

	me, err := h.service.Me(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.log.Warn("profile lookup failed: user not found", "user_id", userID)
			httpx.Error(w, http.StatusNotFound, "user profile not found", httpx.ErrCodeNotFound)
			return
		}

		h.log.Error("internal error fetching user profile", "user_id", userID, "error", err)
		httpx.Error(w, http.StatusInternalServerError, "internal server error", httpx.ErrCodeInternal)
		return
	}

	httpx.JSON(w, http.StatusOK, me)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			httpx.Error(w, http.StatusUnauthorized, "Refresh token cookie is missing", httpx.ErrCodeUnauthorized)
			return
		}

		httpx.Error(w, http.StatusBadRequest, "Invalid request", httpx.ErrBadRequest)
		return
	}

	h.log.Info("Cookie", "val", cookie.Value)

	oldRefreshToken := cookie.Value

	res, err := h.service.Refresh(r.Context(), oldRefreshToken)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrUserNotFound) {
			httpx.Error(w, http.StatusUnauthorized, "Invalid or expired refresh token", httpx.ErrCodeUnauthorized)
			return
		}

		h.log.Error("unexpected error during token refresh handler", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Internal server error", httpx.ErrCodeInternal)
		return
	}

	newCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    res.refreshToken,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   h.env == "production",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	http.SetCookie(w, &newCookie)

	httpx.JSON(w, http.StatusOK, map[string]string{
		"access_token": res.accessToken,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			httpx.Error(w, http.StatusUnauthorized, "Refresh token cookie is missing", httpx.ErrCodeUnauthorized)
			return
		}

		httpx.Error(w, http.StatusBadRequest, "Invalid request", httpx.ErrBadRequest)
		return
	}

	h.log.Info("processing logout", "cookie_found", true)

	rToken := cookie.Value

	err = h.service.Logout(r.Context(), rToken)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			h.log.Warn("logout attempted with non-existent or expired session cookie")
		} else {
			h.log.Error("failed to process logout service layer", "error", err)
			httpx.Error(w, http.StatusInternalServerError, "Internal server error", httpx.ErrCodeInternal)
			return
		}
	}

	newCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   h.env == "production",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}

	http.SetCookie(w, &newCookie)

	httpx.JSON(w, http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}
