package user

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/aryansaud-80/rosterpulse/internal/shared/httpx"
	"github.com/aryansaud-80/rosterpulse/internal/shared/validator"
	goValidator "github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service *UserService
	log     *slog.Logger
}

func NewUserHandler(service *UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		log:     log.With("component", "user.handler"),
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode user request", "error", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid request payload", httpx.ErrCodeInvalidPayload)
		return
	}

	if err := validator.Validate.Struct(&req); err != nil {
		fieldErrors := make(map[string]string)
		if validationErrs, ok := err.(goValidator.ValidationErrors); ok {
			for _, e := range validationErrs {
				field := e.Field()
				switch e.Tag() {
				case "required":
					fieldErrors[field] = field + " is required"
				case "email":
					fieldErrors[field] = "invalid email format"
				case "min":
					fieldErrors[field] = field + " must be at least " + e.Param() + " characters"
				case "oneof":
					fieldErrors[field] = field + " must be one of: " + e.Param()
				default:
					fieldErrors[field] = "invalid value"
				}
			}
		}
		httpx.ValidationError(w, http.StatusBadRequest, "Validation failed", httpx.ErrCodeValidation, fieldErrors)
		return
	}

	res, err := h.service.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			httpx.Error(w, http.StatusConflict, err.Error(), httpx.ErrCodeConflict)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "internal server error", httpx.ErrCodeInternal)
		return
	}

	httpx.JSON(w, http.StatusCreated, res)

}
