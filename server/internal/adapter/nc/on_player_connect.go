package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log/slog"
)

func (p *Producer) OnPlayerConnect(id uint64) {
	p.mxOnPlayerConnect.Lock()
	defer p.mxOnPlayerConnect.Unlock()
	p.onPlayerConnect[id] = struct{}{}
}

func (p *Producer) onPlayerConnectBatch() []*api.Message {
	p.mxOnPlayerConnect.Lock()
	defer p.mxOnPlayerConnect.Unlock()
	defer clear(p.onPlayerConnect)
	batch := make([]*api.Message, 0, len(p.onPlayerConnect))
	for id := range p.onPlayerConnect {
		p.logger.Debug("player has connected",
			slog.Uint64("id", id),
		)
		batch = p.appendToBatch(batch, api.Code_ON_PLAYER_CONNECT, &api.OnPlayerConnect{
			PlayerId: id,
		})
	}
	return batch
}
