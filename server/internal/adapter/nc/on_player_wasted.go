package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log/slog"
)

func (p *Producer) OnPlayerWasted(playerID uint64, playerClass entity.Class, killerID uint64, killerClass entity.Class) {
	p.mxOnPlayerWasted.Lock()
	defer p.mxOnPlayerWasted.Unlock()
	p.onPlayerWasted = append(p.onPlayerWasted, &api.OnPlayerWasted{
		PlayerId:    playerID,
		PlayerClass: api.Class(playerClass),
		KillerId:    killerID,
		KillerClass: api.Class(killerClass),
	})
}

func (p *Producer) onPlayerWastedBatch() []*api.Message {
	p.mxOnPlayerWasted.Lock()
	defer p.mxOnPlayerWasted.Unlock()
	defer func() {
		p.onPlayerWasted = nil
	}()
	batch := make([]*api.Message, 0, len(p.onPlayerWasted))
	for _, msg := range p.onPlayerWasted {
		p.logger.Debug("player died",
			slog.Uint64("id", msg.PlayerId),
			slog.Uint64("killer", msg.KillerId),
		)
		batch = p.appendToBatch(batch, api.Code_ON_PLAYER_WASTED, msg)
	}
	return batch
}
