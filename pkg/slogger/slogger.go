package slogger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"app", "LOCALSEND-AUTO-LOADER",
	)
}
