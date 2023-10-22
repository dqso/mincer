package repository_world

import (
	"context"
	"github.com/dqso/mincer/server/internal/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	logger log.Logger
	pool   *pgxpool.Pool

	projectileIDs chan uint64
}

func NewRepository(ctx context.Context, logger log.Logger, pool *pgxpool.Pool) *Repository {
	r := &Repository{
		logger:        logger.With(log.Module("repo_token")),
		pool:          pool,
		projectileIDs: make(chan uint64, 10),
	}

	go r.acquireProjectileID(ctx)

	return r
}
