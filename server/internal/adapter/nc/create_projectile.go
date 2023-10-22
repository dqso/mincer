package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/internal/log"
	"log/slog"
)

func (p *Producer) CreateProjectile(projectile entity.Projectile) {
	p.mxCreateProjectile.Lock()
	defer p.mxCreateProjectile.Unlock()
	p.createProjectile[projectile.ID()] = projectile
}

func (p *Producer) createProjectileBatch() []*api.Message {
	p.mxCreateProjectile.Lock()
	defer p.mxCreateProjectile.Unlock()
	defer clear(p.createProjectile)
	batch := make([]*api.Message, 0, len(p.createProjectile))
	for _, projectile := range p.createProjectile {
		p.logger.Debug("create projectile",
			slog.Uint64("id", projectile.ID()),
			slog.Uint64("owner", projectile.Owner()),
			log.Damage(projectile.Damage()),
		)
		batch = p.appendToBatch(batch, api.Code_CREATE_PROJECTILE, &api.CreateProjectile{
			Projectile: dtoProjectile(projectile),
		})
	}
	return batch
}
