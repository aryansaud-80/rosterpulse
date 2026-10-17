package app

import (
	"log/slog"
	"net/http"

	"github.com/aryansaud-80/rosterpulse/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Router *http.ServeMux
	db     *pgxpool.Pool
	log    *slog.Logger
}

func NewServer(db *pgxpool.Pool, log *slog.Logger) *Server {
	mux := http.NewServeMux()

	s := &Server{
		Router: mux,
		db:     db,
		log:    log,
	}

	mux.HandleFunc("GET /healthz", handlers.Health)

	return s
}
