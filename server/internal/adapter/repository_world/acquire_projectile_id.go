package repository_world

import (
	"context"
	"log"
	"math"
)

func (r *Repository) AcquireProjectileID() uint64 {
	id, ok := <-r.projectileIDs
	if !ok {
		log.Printf("channel projectileIDs is closed")
		return 0
	}
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
			}
			if err != nil {
				return 0, err
			}
			return id, nil
		}()
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		select {
		case <-ctx.Done():
			return
		case r.projectileIDs <- id:
		}
	}
}
