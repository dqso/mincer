package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log"
)

func (p *Producer) SetPlayerRadius(id uint64, radius float64) {
	p.mxPlayerRadius.Lock()
	defer p.mxPlayerRadius.Unlock()
	p.playerRadius[id] = radius
}

func (p *Producer) setPlayerRadiusBatch() []*api.Message {
	p.mxPlayerRadius.Lock()
	defer p.mxPlayerRadius.Unlock()
	defer clear(p.playerRadius)
	batch := make([]*api.Message, 0, len(p.playerRadius))
	for id, radius := range p.playerRadius {
		msg, err := p.prepareMessage(api.Code_SET_PLAYER_RADIUS, &api.SetPlayerRadius{Id: id, Radius: radius})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
