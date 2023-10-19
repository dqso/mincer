package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) SetPlayerPosition(id uint64, position entity.Point) {
	p.mxPlayerPositions.Lock()
	defer p.mxPlayerPositions.Unlock()
	p.playerPositions[id] = position
}

func (p *Producer) setPlayerPositionBatch() []*api.Message {
	p.mxPlayerPositions.Lock()
	defer p.mxPlayerPositions.Unlock()
	defer clear(p.playerPositions)
	batch := make([]*api.Message, 0, len(p.playerPositions))
	for id, position := range p.playerPositions {
		msg, err := p.prepareMessage(api.Code_SET_PLAYER_POSITION, &api.SetPlayerPosition{Id: id, X: position.X, Y: position.Y})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
