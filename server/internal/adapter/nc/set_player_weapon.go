package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log/slog"
)

func (p *Producer) SetPlayerWeapon(id uint64, weapon entity.Weapon) {
	p.mxPlayerWeapon.Lock()
	defer p.mxPlayerWeapon.Unlock()
	p.playerWeapon[id] = weapon
}

func (p *Producer) setPlayerWeaponBatch() []*api.Message {
	p.mxPlayerWeapon.Lock()
	defer p.mxPlayerWeapon.Unlock()
	defer clear(p.playerWeapon)
	batch := make([]*api.Message, 0, len(p.playerWeapon))
	for id, weapon := range p.playerWeapon {
		p.logger.Debug("set player weapon",
			slog.Uint64("id", id),
			slog.String("name", weapon.Name()),
		)
		batch = p.appendToBatch(batch, api.Code_SET_PLAYER_WEAPON, &api.SetPlayerWeapon{
			Id:     id,
			Weapon: dtoWeapon(weapon),
		})
	}
	return batch
}
