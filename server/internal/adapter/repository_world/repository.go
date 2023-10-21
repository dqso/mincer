package repository_world

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool

	projectileIDs chan uint64
}

func NewRepository(ctx context.Context, pool *pgxpool.Pool) *Repository {
	r := &Repository{
		pool:          pool,
		projectileIDs: make(chan uint64, 10),
	}

	go r.acquireProjectileID(ctx)

	return r
}
