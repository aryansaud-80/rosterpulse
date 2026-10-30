package profile

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aryansaud-80/rosterpulse/internal/auth"
	"github.com/aryansaud-80/rosterpulse/internal/shared/httpx"
	"github.com/aryansaud-80/rosterpulse/internal/shared/validator"
	"github.com/aryansaud-80/rosterpulse/internal/user"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	service IProfileService
	log     *slog.Logger
}

func NewProfileHandler(service IProfileService, log *slog.Logger) *ProfileHandler {
	return &ProfileHandler{
		service: service,
		log:     log.With("component", "profile.handler"),
	}
}

func (h *ProfileHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", httpx.ErrCodeUnauthorized)
		return
	}

	var req CreateProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode profile request", "error", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid request payload", httpx.ErrCodeInvalidPayload)
		return
	}

	if err := validator.Validate.Struct(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Validation failed", httpx.ErrCodeValidation)
		return
	}

	if err := h.service.CreateProfile(r.Context(), userId, req); err != nil {
		if errors.Is(err, ErrProfileAlreadyExists) {
			httpx.Error(w, http.StatusConflict, "Profile already exists for this user", httpx.ErrCodeConflict)
			return
		}
		if errors.Is(err, user.ErrUserNotFound) {
			httpx.Error(w, http.StatusNotFound, "User account not found", httpx.ErrCodeNotFound)
			return
		}

		h.log.Error("failed to create profile", "error", err, "user_id", userId)
		httpx.Error(w, http.StatusInternalServerError, "Internal server error", httpx.ErrCodeInternal)
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]string{
		"message": "profile created successfully",
		"success": "ok",
	})
}

func (h *ProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", httpx.ErrCodeUnauthorized)
		return
	}

	profile, err := h.service.GetProfileByUserID(r.Context(), userId)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			httpx.Error(w, http.StatusNotFound, "Profile not found", httpx.ErrCodeNotFound)
			return
		}

		h.log.Error("failed to fetch user profile", "error", err, "user_id", userId)
		httpx.Error(w, http.StatusInternalServerError, "Internal server error", httpx.ErrCodeInternal)
		return
	}

	httpx.JSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "Unauthorized", httpx.ErrCodeUnauthorized)
		return
	}

	var req UpdateProfileParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode profile request", "error", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid request payload", httpx.ErrCodeInvalidPayload)
		return
	}

	if err := validator.Validate.Struct(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Validation failed", httpx.ErrCodeValidation)
		return
	}

	err = h.service.UpdateProfile(r.Context(), userId, &req)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			httpx.Error(w, http.StatusNotFound, "Profile not found", httpx.ErrCodeNotFound)
			return
		}

		h.log.Error("failed to update user profile", "error", err, "user_id", userId)
		httpx.Error(w, http.StatusInternalServerError, "Internal server error", httpx.ErrCodeInternal)
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]string{
		"message": "profile updated successfully",
		"success": "ok",
	})
}

func (h *ProfileHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	playerId := r.PathValue("playerID")

	if strings.TrimSpace(playerId) == "" {
		httpx.Error(w, http.StatusBadRequest, "Player ID is required", httpx.ErrCodeInvalidPayload)
		return
	}

	if _, err := uuid.Parse(playerId); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid player ID format", httpx.ErrCodeInvalidPayload)
		return
	}

	profile, err := h.service.GetPublicProfileByID(r.Context(), playerId)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			httpx.Error(w, http.StatusNotFound, "Profile not found", httpx.ErrCodeNotFound)
			return
		}

		h.log.Error("failed to fetch public profile", "error", err, "player_id", playerId)
		httpx.Error(w, http.StatusInternalServerError, "Internal server error", httpx.ErrCodeInternal)
		return
	}

	httpx.JSON(w, http.StatusOK, profile)
}
