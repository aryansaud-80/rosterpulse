package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aryansaud-80/rosterpulse/internal/config"
	"github.com/aryansaud-80/rosterpulse/internal/db"
	"github.com/aryansaud-80/rosterpulse/internal/handlers"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx, cfg.DATABASE_URL)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	log.Printf("connected to database")

	defer pool.Close()

	log.Printf("server starting ....")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("server is listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("could not listen on %s due to %s", srv.Addr, err.Error())
	}

}
