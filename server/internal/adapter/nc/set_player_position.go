package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/internal/log"
	"log/slog"
)

func (p *Producer) SetPlayerPosition(id uint64, position entity.Point) {
	p.mxPlayerPositions.Lock()
	defer p.mxPlayerPositions.Unlock()
	p.playerPositions[id] = position
}

func (p *Producer) setPlayerPositionBatch() []*api.Message {
	p.mxPlayerPositions.Lock()
	defer p.mxPlayerPositions.Unlock()
	defer clear(p.playerPositions)
	batch := make([]*api.Message, 0, len(p.playerPositions))
	for id, position := range p.playerPositions {
		p.logger.Debug("set player position",
			slog.Uint64("id", id),
			log.Point(position),
		)
		batch = p.appendToBatch(batch, api.Code_SET_PLAYER_POSITION, &api.SetPlayerPosition{
			Id: id,
			X:  position.X, // TODO api.Point
			Y:  position.Y,
		})
	}
	return batch
}
