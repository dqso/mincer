package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log/slog"
)

func (p *Producer) OnPlayerAttacked(id uint64, directionAim float64) {
	p.mxOnPlayerAttacked.Lock()
	defer p.mxOnPlayerAttacked.Unlock()
	p.onPlayerAttacked = append(p.onPlayerAttacked, &api.OnPlayerAttacked{
		Id:           id,
		DirectionAim: directionAim,
	})
}

func (p *Producer) onPlayerAttackedBatch() []*api.Message {
	p.mxOnPlayerAttacked.Lock()
	defer p.mxOnPlayerAttacked.Unlock()
	defer func() {
		p.onPlayerAttacked = nil
	}()
	batch := make([]*api.Message, 0, len(p.onPlayerAttacked))
	for _, msg := range p.onPlayerAttacked {
		p.logger.Debug("player attacked",
			slog.Uint64("id", msg.Id),
			slog.Float64("direction_aim", msg.DirectionAim),
		)
		batch = p.appendToBatch(batch, api.Code_ON_PLAYER_ATTACKED, msg)
	}
	return batch
}
