package repository_token

import (
	"context"
	"log/slog"
)

func (r Repository) AcquireClientID(ctx context.Context) (uint64, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var clientID uint64
	const stmt = `SELECT nextval('client_id')`
	err = conn.QueryRow(ctx, stmt).Scan(&clientID)
	if err != nil {
		return 0, err
	}

	r.logger.Debug("client id has been acquired",
		slog.Uint64("id", clientID),
	)
	return clientID, nil
}
