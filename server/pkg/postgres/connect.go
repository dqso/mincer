package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type config interface {
	PostgresHost() string
	PostgresPort() string
	PostgresDatabase() string
	PostgresUsername() string
	PostgresPassword() string
}

func Connect(ctx context.Context, config config) (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.PostgresUsername(), config.PostgresPassword(),
		config.PostgresHost(), config.PostgresPort(), config.PostgresDatabase(),
	))
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
