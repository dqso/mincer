package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) SetPlayerClass(id uint64, class entity.Class) {
	p.mxPlayerClasses.Lock()
	defer p.mxPlayerClasses.Unlock()
	p.playerClasses[id] = class
}

func (p *Producer) setPlayerClassBatch() []*api.Message {
	p.mxPlayerClasses.Lock()
	defer p.mxPlayerClasses.Unlock()
	defer clear(p.playerClasses)
	batch := make([]*api.Message, 0, len(p.playerClasses))
	for id, class := range p.playerClasses {
		msg, err := p.prepareMessage(api.Code_SET_PLAYER_CLASS, &api.SetPlayerClass{Id: id, Class: api.Class(class)})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
