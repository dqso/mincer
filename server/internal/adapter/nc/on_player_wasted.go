package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log/slog"
)

func (p *Producer) OnPlayerWasted(id uint64, killer uint64) {
	p.mxOnPlayerWasted.Lock()
	defer p.mxOnPlayerWasted.Unlock()
	p.onPlayerWasted = append(p.onPlayerWasted, &api.OnPlayerWasted{
		Id:     id,
		Killer: killer,
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
			slog.Uint64("id", msg.Id),
			slog.Uint64("killer", msg.Killer),
		)
		batch = p.appendToBatch(batch, api.Code_ON_PLAYER_WASTED, msg)
	}
	return batch
}
