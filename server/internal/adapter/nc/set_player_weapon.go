package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) SetPlayerWeapon(id uint64, w entity.Weapon) {
	p.mxPlayerWeapon.Lock()
	defer p.mxPlayerWeapon.Unlock()
	p.playerWeapon[id] = dtoWeapon(w)
}

func (p *Producer) setPlayerWeaponBatch() []*api.Message {
	p.mxPlayerWeapon.Lock()
	defer p.mxPlayerWeapon.Unlock()
	defer clear(p.playerWeapon)
	batch := make([]*api.Message, 0, len(p.playerWeapon))
	for id, weapon := range p.playerWeapon {
		msg, err := p.prepareMessage(api.Code_SET_PLAYER_WEAPON, &api.SetPlayerWeapon{Id: id, Weapon: weapon})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
