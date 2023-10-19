package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log"
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
	for playerID := range p.onPlayerDisconnect {
		msg, err := p.prepareMessage(api.Code_ON_PLAYER_DISCONNECT, &api.OnPlayerDisconnect{
			PlayerId: playerID,
		})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
