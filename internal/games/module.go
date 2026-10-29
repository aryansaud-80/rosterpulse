package games

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterGameModule(
	mux *http.ServeMux,
	db *pgxpool.Pool,
	log *slog.Logger,
) {
	gameRepo := NewGameRepository(db, log)
	gameService := NewGameService(gameRepo, log)
	gameHandler := NewGameHandler(gameService, log)

	mux.HandleFunc("GET /api/games", gameHandler.GetGames)
	mux.HandleFunc("GET /api/games/{gameid}", gameHandler.GetGameById)
}
