package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Env          string
	DATABASE_URL string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")

	if port == "" {
		panic("PORT is required")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL is required")
	}

	return Config{
		Port:         port,
		Env:          env,
		DATABASE_URL: databaseURL,
	}

}
