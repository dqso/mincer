package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log"
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
	for playerID := range p.onPlayerConnect {
		msg, err := p.prepareMessage(api.Code_ON_PLAYER_CONNECT, &api.OnPlayerConnect{
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
