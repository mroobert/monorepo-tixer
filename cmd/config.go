package main

import (
	"fmt"
	"time"

	"github.com/mroobert/monorepo-tixer/env"
	"github.com/mroobert/monorepo-tixer/httpio"
	"github.com/mroobert/monorepo-tixer/psql"
)

// Config represents the application configuration details.
type config struct {
	Env      string // the environment the application is running in.
	Server   httpio.ServerConfig
	Database psql.DbConfig
}

// NewConfig creates a new instance of Config.
// It loads the configuration details from the environment variables.
func newConfig() (*config, error) {

	environment, err := env.LoadEnv("ENVIRONMENT")
	if err != nil {
		return nil, fmt.Errorf("loading ENVIRONMENT failed: %w", err)
	}

	// Load the server configuration.
	serverAddr := env.LoadEnvOrDefault("SERVER_ADDR", "localhost:8080")

	serverIdleTimeout, err := env.LoadDurationEnvOrDefault("SERVER_IDLE_TIMEOUT", 120*time.Second)
	if err != nil {
		return nil, fmt.Errorf("loading SERVER_IDLE_TIMEOUT failed: %w", err)
	}

	serverReadTimeout, err := env.LoadDurationEnvOrDefault("SERVER_READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("loading SERVER_READ_TIMEOUT failed: %w", err)
	}

	serverShutdownTimeout, err := env.LoadDurationEnvOrDefault("SERVER_SHUTDOWN_TIMEOUT", 20*time.Second)
	if err != nil {
		return nil, fmt.Errorf("loading SERVER_SHUTDOWN_TIMEOUT failed: %w", err)
	}

	serverWriteTimeout, err := env.LoadDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("loading SERVER_WRITE_TIMEOUT failed: %w", err)
	}

	serverMaxReqBodySize, err := env.LoadInt32EnvOrDefault("SERVER_MAX_REQ_BODY_SIZE", 1048576)
	if err != nil {
		return nil, fmt.Errorf("loading SERVER_MAX_REQ_BODY_SIZE failed: %w", err)
	}

	serverConfig := httpio.ServerConfig{
		Addr:            serverAddr,
		IdleTimeout:     serverIdleTimeout,
		ReadTimeout:     serverReadTimeout,
		ShutdownTimeout: serverShutdownTimeout,
		WriteTimeout:    serverWriteTimeout,
		MaxReqBodySize:  serverMaxReqBodySize,
	}

	// Load the database configuration.
	dbDSN, err := env.LoadEnv("DB_DSN")
	if err != nil {
		return nil, fmt.Errorf("loading DB_DSN failed: %w", err)
	}

	dbMaxOpenConns, err := env.LoadInt32EnvOrDefault("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return nil, fmt.Errorf("loading DB_MAX_OPEN_CONNS failed: %w", err)
	}

	dbMinConns, err := env.LoadInt32EnvOrDefault("DB_MIN_CONNS", 5)
	if err != nil {
		return nil, fmt.Errorf("loading DB_MIN_CONNS failed: %w", err)
	}

	dbMaxConnIdleTime, err := env.LoadDurationEnvOrDefault("DB_MAX_CONN_IDLE_TIME", 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("loading DB_MAX_CONN_IDLE_TIME failed: %w", err)
	}

	dbQueryTimeout, err := env.LoadDurationEnvOrDefault("DB_QUERY_TIMEOUT", 3*time.Second)
	if err != nil {
		return nil, fmt.Errorf("loading DB_QUERY_TIMEOUT failed: %w", err)
	}

	dbConfig := psql.DbConfig{
		DSN:             dbDSN,
		MaxOpenConns:    dbMaxOpenConns,
		MinConns:        dbMinConns,
		MaxConnIdleTime: dbMaxConnIdleTime,
		QueryTimeout:    dbQueryTimeout,
	}

	return &config{
		Env:      environment,
		Server:   serverConfig,
		Database: dbConfig,
	}, nil
}
