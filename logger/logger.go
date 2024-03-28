// This package provides support for configuring the log system.
package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mroobert/monorepo-tixer/httpio/rcontext"
)

// SetLogger sets the logger based on the environment.
func SetLogger(env string) {
	var logger *slog.Logger

	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	switch env {
	case "local":
		logger = slog.New(&localHandler{
			textHandler: slog.NewTextHandler(os.Stdout, handlerOptions),
		})
	case "development", "staging", "production":
		logger = slog.New(&nonLocalHandler{
			jsonHandler: slog.NewJSONHandler(os.Stdout, handlerOptions),
		})
	}

	slog.SetDefault(logger)
}

// localHandler is a custom handler that writes to os.Stdout and adds a newline.
// It also adds the request_id to the log record if it is present in the context.
type localHandler struct {
	textHandler slog.Handler
}

func (h *localHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.textHandler.Enabled(ctx, lvl)
}

func (h *localHandler) Handle(ctx context.Context, r slog.Record) error {
	requestInfo := rcontext.GetRequestInfo(ctx)
	if requestInfo != nil {
		r.AddAttrs(slog.String("request_id", requestInfo.RequestID))
	}

	err := h.textHandler.Handle(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to write log record: %w", err)
	}
	_, err = os.Stdout.Write([]byte("\n"))
	return fmt.Errorf("failed to write newline for log: %w", err)
}

func (h *localHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &localHandler{textHandler: h.textHandler.WithAttrs(attrs)}
}

func (h *localHandler) WithGroup(name string) slog.Handler {
	return &localHandler{textHandler: h.textHandler.WithGroup(name)}
}

// nonLocalHandler is a custom handler that writes
type nonLocalHandler struct {
	jsonHandler slog.Handler
}

func (h *nonLocalHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.jsonHandler.Enabled(ctx, lvl)
}

func (h *nonLocalHandler) Handle(ctx context.Context, r slog.Record) error {
	requestInfo := rcontext.GetRequestInfo(ctx)
	if requestInfo != nil {
		r.AddAttrs(slog.String("request_id", requestInfo.RequestID))
	}

	err := h.jsonHandler.Handle(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to write log record: %w", err)
	}

	return nil
}

func (h *nonLocalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &localHandler{textHandler: h.jsonHandler.WithAttrs(attrs)}
}

func (h *nonLocalHandler) WithGroup(name string) slog.Handler {
	return &localHandler{textHandler: h.jsonHandler.WithGroup(name)}
}
