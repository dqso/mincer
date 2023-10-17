package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) OnPlayerChange(player entity.Player) {
	p.mxOnPlayerChange.Lock()
	defer p.mxOnPlayerChange.Unlock()
	p.onPlayerChange[player.ID()] = player
}

func (p *Producer) onPlayerChangeBatch() []*api.Message {
	p.mxOnPlayerChange.Lock()
	defer p.mxOnPlayerChange.Unlock()
	batch := make([]*api.Message, 0, len(p.onPlayerChange))
	for _, player := range p.onPlayerChange {
		msg, err := p.prepareMessage(api.Code_ON_PLAYER_CHANGE, &api.OnPlayerChange{Player: dtoPlayerToPublicPlayer(player)})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	clear(p.onPlayerChange)
	return batch
}
