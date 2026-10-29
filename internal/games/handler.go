package games

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/aryansaud-80/rosterpulse/internal/shared/httpx"
)

type GameHandler struct {
	service IGameService
	log     *slog.Logger
}

func NewGameHandler(service IGameService, log *slog.Logger) *GameHandler {
	return &GameHandler{
		service: service,
		log:     log.With("component", "games.handler"),
	}
}

func (h *GameHandler) GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.service.GetGames(r.Context())
	if err != nil {
		h.log.Error("failed to fetch games", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Internal Server Error", httpx.ErrCodeInternal)
		return
	}

	httpx.JSON(w, http.StatusOK, games)
}

func (h *GameHandler) GetGameById(w http.ResponseWriter, r *http.Request) {
	gameid := r.PathValue("gameid")

	game, err := h.service.GetGameById(r.Context(), gameid)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidGameID):
			httpx.Error(w, http.StatusBadRequest, "Invalid game ID format", httpx.ErrBadRequest)
			return
		case errors.Is(err, ErrGameNotFound):
			httpx.Error(w, http.StatusNotFound, "Game not found", httpx.ErrCodeNotFound)
			return
		default:
			h.log.Error("failed to fetch game by id", "game_id", gameid, "error", err)
			httpx.Error(w, http.StatusInternalServerError, "Internal Server Error", httpx.ErrCodeInternal)
			return
		}
	}

	httpx.JSON(w, http.StatusOK, game)
}
