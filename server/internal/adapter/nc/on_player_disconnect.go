package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log/slog"
)

func (p *Producer) OnPlayerDisconnect(id uint64) {
	p.mxOnPlayerDisconnect.Lock()
	defer p.mxOnPlayerDisconnect.Unlock()
	p.onPlayerDisconnect[id] = struct{}{}
}

func (p *Producer) onPlayerDisconnectBatch() []*api.Message {
	p.mxOnPlayerDisconnect.Lock()
	defer p.mxOnPlayerDisconnect.Unlock()
	defer clear(p.onPlayerDisconnect)
	batch := make([]*api.Message, 0, len(p.onPlayerDisconnect))
	for id := range p.onPlayerDisconnect {
		p.logger.Debug("player has disconnected",
			slog.Uint64("id", id),
		)
		batch = p.appendToBatch(batch, api.Code_ON_PLAYER_DISCONNECT, &api.OnPlayerDisconnect{
			PlayerId: id,
		})
	}
	return batch
}
