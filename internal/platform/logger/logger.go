package logger

import (
	"log/slog"
	"os"

	"github.com/aryansaud-80/rosterpulse/internal/config"
)

func NewLogger() *slog.Logger {
	env := config.MustLoad().Env
	var level slog.Level
	addSource := false

	switch env {
	case "production":
		level = slog.LevelInfo
		addSource = false
	case "development":
		level = slog.LevelDebug
		addSource = true
	default:
		level = slog.LevelDebug
		addSource = true
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: addSource,
		Level:     level,
	})

	return slog.New(handler)
}
