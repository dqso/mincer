package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log"
)

func (p *Producer) SetPlayerHP(id uint64, hp int32) {
	p.mxPlayerHP.Lock()
	defer p.mxPlayerHP.Unlock()
	p.playerHP[id] = hp
}

func (p *Producer) setPlayerHPBatch() []*api.Message {
	p.mxPlayerHP.Lock()
	defer p.mxPlayerHP.Unlock()
	defer clear(p.playerHP)
	batch := make([]*api.Message, 0, len(p.playerHP))
	for id, hp := range p.playerHP {
		msg, err := p.prepareMessage(api.Code_SET_PLAYER_HP, &api.SetPlayerHP{Id: id, Hp: hp})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
