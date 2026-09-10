package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds database connection parameters.
type Config struct {
	URL          string
	MaxConns     int32
	MinConns     int32
	MaxIdleTime  time.Duration
	MaxConnLife  time.Duration
	DialTimeout  time.Duration
}

// DefaultConfig returns standard connection settings for PostgreSQL.
func DefaultConfig() Config {
	return Config{
		URL:         "postgres://postgres:postgres@localhost:5432/hydravms?sslmode=disable",
		MaxConns:    20,
		MinConns:    2,
		MaxIdleTime: 5 * time.Minute,
		MaxConnLife: 30 * time.Minute,
		DialTimeout: 3 * time.Second,
	}
}

// NewPool initializes a managed PostgreSQL connection pool.
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid postgres connection string: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnIdleTime = cfg.MaxIdleTime
	poolCfg.MaxConnLifetime = cfg.MaxConnLife

	// Hardening runtime session parameters
	if poolCfg.ConnConfig.RuntimeParams == nil {
		poolCfg.ConnConfig.RuntimeParams = make(map[string]string)
	}
	poolCfg.ConnConfig.RuntimeParams["application_name"] = "hydravms-control-plane"
	poolCfg.ConnConfig.RuntimeParams["statement_timeout"] = "5000"                      // 5s limit
	poolCfg.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] = "10000" // 10s limit

	ctxTimeout, cancel := context.WithTimeout(ctx, cfg.DialTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctxTimeout, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctxTimeout); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return pool, nil
}
