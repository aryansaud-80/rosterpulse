package user

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserModule(mux *http.ServeMux, db *pgxpool.Pool, log *slog.Logger) Iuser {

	userRepo := NewUserRepository(db, log)
	userService := NewUserService(userRepo, log)
	userHandler := NewUserHandler(userService, log)

	mux.HandleFunc("POST /api/users", userHandler.Register)

	return userRepo
}
