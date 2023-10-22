package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/internal/log"
	"log/slog"
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
		p.logger.Debug("set projectile position",
			slog.Uint64("id", id),
			log.Point(position),
		)
		batch = p.appendToBatch(batch, api.Code_SET_PROJECTILE_POSITION, &api.SetProjectilePosition{
			Id:       id,
			Position: dtoPoint(position),
		})
	}
	return batch
}
