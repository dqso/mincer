package repository_token

import "context"

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

	return clientID, nil
}
