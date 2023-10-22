package repository_world

import (
	"context"
	"log/slog"
	"math"
	"math/rand"
)

func (r *Repository) AcquireProjectileID() uint64 {
	id, ok := <-r.projectileIDs
	if !ok {
		r.logger.Warn("unable to acquire the projectile id because projectileIDs channel is closed. Random numbers are used")
		return rand.Uint64()
	}
	r.logger.Debug("projectile id has been acquired",
		slog.Uint64("id", id),
	)
	return id
}

func (r *Repository) acquireProjectileID(ctx context.Context) {
	defer close(r.projectileIDs)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		id, err := func() (uint64, error) {
			conn, err := r.pool.Acquire(ctx)
			if err != nil {
				return 0, err
			}
			defer conn.Release()
			var id uint64
			const stmt = `SELECT nextval('projectile_id')`
			err = conn.QueryRow(ctx, stmt).Scan(&id)
			if id == math.MaxInt64 || err != nil {
				const stmt = `ALTER SEQUENCE projectile_id RESTART 1;`
				if _, err := conn.Exec(ctx, stmt); err != nil {
					return 0, err
				}
				r.logger.Info("projectile id pool has been reset")
			}
			if err != nil {
				return 0, err
			}
			return id, nil
		}()
		if err != nil {
			r.logger.Error("unable to acquire the projectile id")
			continue
		}
		select {
		case <-ctx.Done():
			return
		case r.projectileIDs <- id:
		}
	}
}
