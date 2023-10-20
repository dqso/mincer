package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log"
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
	for _, wasted := range p.onPlayerWasted {
		msg, err := p.prepareMessage(api.Code_ON_PLAYER_WASTED, wasted)
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
