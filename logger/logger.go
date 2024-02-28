// This package provides support for configuring the log system.
package logger

import (
	"log/slog"
	"os"
)

// SetLogger sets the logger based on the environment.
func SetLogger(env string) {
	var logger *slog.Logger

	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	case "development", "staging", "production":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	}

	slog.SetDefault(logger)
}
