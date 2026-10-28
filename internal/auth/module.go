package auth

import (
	"log/slog"
	"net/http"

	"github.com/aryansaud-80/rosterpulse/internal/config"
	"github.com/aryansaud-80/rosterpulse/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthModule(
	mux *http.ServeMux,
	db *pgxpool.Pool,
	userRepo user.Iuser,
	log *slog.Logger,
	cfg *config.Config,
	tm *TokenManager,
	authMw func(http.Handler) http.Handler,
) {
	sessionRepo := NewSessionRepository(db, log)
	authService := NewAuthService(sessionRepo, userRepo, log, tm)
	authHandler := NewAuthHandler(authService, log, cfg.Env)

	mux.HandleFunc("POST /api/auth/refresh", authHandler.Refresh)

	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("GET /api/auth/me", authMw(http.HandlerFunc(authHandler.Me)))
	mux.Handle("POST /api/auth/logout", authMw(http.HandlerFunc(authHandler.Logout)))
}
