package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) CreateProjectile(projectile entity.Projectile) {
	p.mxCreateProjectile.Lock()
	defer p.mxCreateProjectile.Unlock()
	p.createProjectile[projectile.ID()] = &api.CreateProjectile{Projectile: dtoProjectile(projectile)}
}

func (p *Producer) createProjectileBatch() []*api.Message {
	p.mxCreateProjectile.Lock()
	defer p.mxCreateProjectile.Unlock()
	defer clear(p.createProjectile)
	batch := make([]*api.Message, 0, len(p.createProjectile))
	for _, projectile := range p.createProjectile {
		msg, err := p.prepareMessage(api.Code_CREATE_PROJECTILE, projectile)
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
