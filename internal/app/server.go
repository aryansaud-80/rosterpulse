package app

import (
	"log/slog"
	"net/http"

	"github.com/aryansaud-80/rosterpulse/internal/auth"
	"github.com/aryansaud-80/rosterpulse/internal/config"
	"github.com/aryansaud-80/rosterpulse/internal/handlers"
	"github.com/aryansaud-80/rosterpulse/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Router *http.ServeMux
	db     *pgxpool.Pool
	log    *slog.Logger
	cfg    *config.Config
}

func NewServer(db *pgxpool.Pool, log *slog.Logger, cfg *config.Config) *Server {
	mux := http.NewServeMux()

	s := &Server{
		Router: mux,
		db:     db,
		log:    log,
		cfg:    cfg,
	}

	mux.HandleFunc("GET /healthz", handlers.Health)

	userRepo := user.RegisterUserModule(mux, db, log)
	auth.RegisterAuthModule(mux, db, userRepo, log, cfg)

	return s
}
