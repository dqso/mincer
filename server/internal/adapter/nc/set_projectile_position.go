package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) SetProjectilePosition(id uint64, position entity.Point) {
	p.mxProjectilePositions.Lock()
	defer p.mxProjectilePositions.Unlock()
	p.projectilePositions[id] = position
}

func (p *Producer) setProjectilePositionBatch() []*api.Message {
	p.mxProjectilePositions.Lock()
	defer p.mxProjectilePositions.Unlock()
	defer clear(p.projectilePositions)
	batch := make([]*api.Message, 0, len(p.projectilePositions))
	for id, position := range p.projectilePositions {
		msg, err := p.prepareMessage(api.Code_SET_PROJECTILE_POSITION, &api.SetProjectilePosition{Id: id, Position: dtoPoint(position)})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
