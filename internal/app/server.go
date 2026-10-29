package app

import (
	"log/slog"
	"net/http"

	"github.com/aryansaud-80/rosterpulse/internal/auth"
	"github.com/aryansaud-80/rosterpulse/internal/config"
	"github.com/aryansaud-80/rosterpulse/internal/games"
	"github.com/aryansaud-80/rosterpulse/internal/handlers"
	"github.com/aryansaud-80/rosterpulse/internal/middleware"
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

	tokenManager := auth.NewTokenManager(
		cfg.JWT.SecretKey,
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)

	authMw := middleware.AuthMiddleware(tokenManager)

	auth.RegisterAuthModule(mux, db, userRepo, log, cfg, tokenManager, authMw)

	games.RegisterGameModule(mux, db, log)

	return s
}
