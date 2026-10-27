package auth

import (
	"log/slog"
	"net/http"

	"github.com/aryansaud-80/rosterpulse/internal/config"
	"github.com/aryansaud-80/rosterpulse/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthModule(mux *http.ServeMux, db *pgxpool.Pool, userRepo user.Iuser, log *slog.Logger, cfg *config.Config) {
	tokenManager := NewTokenManager(
		cfg.JWT.SecretKey,
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)

	sessionRepo := NewSessionRepository(db, log)
	authService := NewAuthService(sessionRepo, userRepo, log, tokenManager)
	authHandler := NewAuthHandler(authService, log, cfg.Env)

	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
}
