package handlers

import (
	"net/http"

	"github.com/aryansaud-80/rosterpulse/internal/shared/response"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response.JSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}
