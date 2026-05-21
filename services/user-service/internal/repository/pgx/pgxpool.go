package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/artlink52/ecommerce_backend/services/user-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewPool(ctx context.Context, cfg config.PostgresConfig) (*Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping pgxpool: %w", err)
	}

	return &Pool{
		Pool:      pool,
		opTimeout: cfg.Timeout,
	}, nil
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}
