package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mroobert/monorepo-tixer/env"
	"github.com/mroobert/monorepo-tixer/httpio"
	"github.com/mroobert/monorepo-tixer/logger"
	"github.com/mroobert/monorepo-tixer/psql"
)

func main() {
	err := env.LoadEnvFile()
	if err != nil {
		fmt.Println("failed to load the environment file: ", err)
		os.Exit(1)
	}

	ctx := context.Background()
	app, err := NewApplication(ctx)
	if err != nil {
		fmt.Println("failed to create the application: ", err)
		os.Exit(1)
	}
	logger.SetLogger(app.Config.Env)

	if err := app.Run(ctx); err != nil {
		slog.Error("failed to run the application", err)
		os.Exit(1)
	}
}

// Application holds the dependencies for the web application.
type Application struct {
	Config *config
	Server *httpio.Server
}

// NewApplication creates a new configured Application.
func NewApplication(ctx context.Context) (*Application, error) {
	cfg, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config failed: %w", err)
	}

	dbPool, err := psql.NewPool(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("connecting to db failed: %w", err)
	}

	server := httpio.NewServer(cfg.Server, cfg.Env)
	server.TicketRepository = psql.NewTicketRepository(dbPool, cfg.Database.QueryTimeout)

	return &Application{
		Config: cfg,
		Server: server,
	}, nil
}

// Run performs the startup sequence.
func (a *Application) Run(ctx context.Context) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("starting the server",
			slog.String("addr", a.Config.Server.Addr),
			slog.String("env", a.Config.Env),
		)

		serverErrors <- a.Server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		slog.Info("shutdown signal received", slog.String("signal", sig.String()))
		defer slog.Info("shutdown complete", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(ctx, a.Config.Server.ShutdownTimeout)
		defer cancel()

		if err := a.Server.Shutdown(ctx); err != nil {
			a.Server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
