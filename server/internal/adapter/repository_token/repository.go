package repository_token

import (
	"github.com/dqso/mincer/server/internal/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	logger log.Logger
	pool   *pgxpool.Pool
}

func NewRepository(logger log.Logger, pool *pgxpool.Pool) *Repository {
	return &Repository{
		logger: logger.With(log.Module("repo_token")),
		pool:   pool,
	}
}
