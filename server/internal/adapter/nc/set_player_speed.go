package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log"
)

func (p *Producer) SetPlayerSpeed(id uint64, speed float64) {
	p.mxPlayerSpeed.Lock()
	defer p.mxPlayerSpeed.Unlock()
	p.playerSpeed[id] = speed
}

func (p *Producer) setPlayerSpeedBatch() []*api.Message {
	p.mxPlayerSpeed.Lock()
	defer p.mxPlayerSpeed.Unlock()
	defer clear(p.playerSpeed)
	batch := make([]*api.Message, 0, len(p.playerSpeed))
	for id, speed := range p.playerSpeed {
		msg, err := p.prepareMessage(api.Code_SET_PLAYER_SPEED, &api.SetPlayerSpeed{Id: id, Speed: speed})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
