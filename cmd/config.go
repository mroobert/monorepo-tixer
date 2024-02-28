package main

import (
	"fmt"
	"time"

	"github.com/mroobert/monorepo-tixer/env"
	"github.com/mroobert/monorepo-tixer/http"
)

// Config represents the application configuration details.
type Config struct {
	Environment string
	Server      http.ServerConfig
}

func LoadConfig() (Config, error) {

	environment, err := env.LoadEnv("ENVIRONMENT")
	if err != nil {
		return Config{}, fmt.Errorf("loading ENVIRONMENT failed: %w", err)
	}

	// Load the server configuration.
	serverAddr := env.LoadEnvOrDefault("SERVER_ADDR", "localhost:8080")

	serverIdleTimeout, err := env.LoadDurationEnvOrDefault("SERVER_IDLE_TIMEOUT", 120*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("loading SERVER_IDLE_TIMEOUT failed: %w", err)
	}

	serverReadTimeout, err := env.LoadDurationEnvOrDefault("SERVER_READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("loading SERVER_READ_TIMEOUT failed: %w", err)
	}

	serverShutdownTimeout, err := env.LoadDurationEnvOrDefault("SERVER_SHUTDOWN_TIMEOUT", 20*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("loading SERVER_SHUTDOWN_TIMEOUT failed: %w", err)
	}

	serverWriteTimeout, err := env.LoadDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, fmt.Errorf("loading SERVER_WRITE_TIMEOUT failed: %w", err)
	}

	serverConfig := http.ServerConfig{
		Addr:            serverAddr,
		IdleTimeout:     serverIdleTimeout,
		ReadTimeout:     serverReadTimeout,
		ShutdownTimeout: serverShutdownTimeout,
		WriteTimeout:    serverWriteTimeout,
	}

	return Config{
		Environment: environment,
		Server:      serverConfig,
	}, nil
}
