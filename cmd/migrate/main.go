package main

import (
	"log"
	"os"

	"github.com/aryansaud-80/rosterpulse/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: migrate <up | down>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations",
		cfg.DATABASE_URL,
	)

	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("could not run migrations up: %v", err)
		}
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("could not run migrations down: %v", err)
		}
	default:
		log.Fatalf("invalid argument: %s", os.Args[1])
	}
}
