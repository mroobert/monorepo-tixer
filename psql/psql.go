// The psql package provides support for working with a PostgreSQL database.
package psql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DbConfig represents the configuration details for the database connection.
type DbConfig struct {
	DSN             string        // data source name
	MaxOpenConns    int32         // limit on the number of ‘open’ connections (in-use + idle connections)
	MinConns        int32         // minimum size of the pool
	MaxConnIdleTime time.Duration // sets the maximum length of time that a connection can be idle for before it is marked as expired
	QueryTimeout    time.Duration // sets the maximum time a query can run before it is canceled
}

// NewPool creates a new connection pool to the database.
func NewPool(cfg DbConfig) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}
	dbConfig.MaxConns = cfg.MaxOpenConns
	dbConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	dbConfig.MinConns = cfg.MinConns

	pool, err := pgxpool.NewWithConfig(context.TODO(), dbConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
