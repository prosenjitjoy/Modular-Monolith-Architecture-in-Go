package logger

import (
	"log/slog"
	"os"
)

func New(logLevel string) *slog.Logger {
	option := &slog.HandlerOptions{
		// AddSource: true,
		Level: toSlogLevel(logLevel),
	}

	return slog.New(slog.NewTextHandler(os.Stdout, option))
}

func toSlogLevel(logLevel string) slog.Leveler {
	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
