package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aryansaud-80/rosterpulse/internal/app"
	"github.com/aryansaud-80/rosterpulse/internal/config"
	"github.com/aryansaud-80/rosterpulse/internal/middleware"
	"github.com/aryansaud-80/rosterpulse/internal/platform/db"
	"github.com/aryansaud-80/rosterpulse/internal/platform/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx, cfg.DATABASE_URL)
	if err != nil {
		log.Error("could not connect to database", "error", err)
		os.Exit(1)
	}

	log.Info("connected to database")

	defer pool.Close()

	log.Info("server starting ....")

	srvContainer := app.NewServer(pool, log, &cfg)

	loggerHandler := middleware.RequestLogger(log)(srvContainer.Router)
	finalHandler := middleware.RequestId(loggerHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      finalHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Info("server is listening", "address", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Error("could not listen on server", "addr", srv.Addr, "error", err)
		os.Exit(1)
	}
}
