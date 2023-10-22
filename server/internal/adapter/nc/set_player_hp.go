package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"log/slog"
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
		p.logger.Debug("set player hp",
			slog.Uint64("id", id),
			slog.Int64("hp", int64(hp)),
		)
		batch = p.appendToBatch(batch, api.Code_SET_PLAYER_HP, &api.SetPlayerHP{
			Id: id,
			Hp: hp,
		})
	}
	return batch
}
