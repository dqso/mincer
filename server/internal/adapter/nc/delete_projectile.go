package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log/slog"
)

func (p *Producer) DeleteProjectile(id uint64) {
	p.mxDeleteProjectile.Lock()
	defer p.mxDeleteProjectile.Unlock()
	p.deleteProjectile[id] = struct{}{}
}

func (p *Producer) deleteProjectileBatch() []*api.Message {
	p.mxDeleteProjectile.Lock()
	defer p.mxDeleteProjectile.Unlock()
	defer clear(p.deleteProjectile)
	batch := make([]*api.Message, 0, len(p.deleteProjectile))
	for id := range p.deleteProjectile {
		p.logger.Debug("delete projectile",
			slog.Uint64("id", id),
		)
		batch = p.appendToBatch(batch, api.Code_DELETE_PROJECTILE, &api.DeleteProjectile{
			Id: id,
		})
	}
	return batch
}
