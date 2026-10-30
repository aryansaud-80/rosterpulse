package profile

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterPlayerModule(
	mux *http.ServeMux,
	db *pgxpool.Pool,
	log *slog.Logger,
	authMw func(http.Handler) http.Handler,
) {
	profileRepo := NewProfileRepository(db, log)
	profileService := NewProfileService(profileRepo, log)
	profileHandler := NewProfileHandler(profileService, log)

	mux.Handle("POST /api/player/profile", authMw(http.HandlerFunc(profileHandler.Create)))
	mux.Handle("GET /api/player/profile/me", authMw(http.HandlerFunc(profileHandler.GetMyProfile)))
	mux.Handle("PATCH /api/player/profile/me", authMw(http.HandlerFunc(profileHandler.UpdateMyProfile)))

	mux.HandleFunc("GET /api/players/{playerID}", profileHandler.GetPublicProfile)

}
