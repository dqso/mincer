package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) SetPlayerStats(id uint64, class entity.Class, radius, speed float64, maxHP int64, maxCoolDown, power float64) {
	p.mxPlayerStats.Lock()
	defer p.mxPlayerStats.Unlock()
	p.playerStats[id] = &api.PlayerStats{
		Class:       api.Class(class),
		Radius:      radius,
		Speed:       speed,
		MaxHP:       maxHP,
		MaxCoolDown: maxCoolDown,
		Power:       power,
	}
}

func (p *Producer) setPlayerStatsBatch() []*api.Message {
	p.mxPlayerStats.Lock()
	defer p.mxPlayerStats.Unlock()
	defer clear(p.playerStats)
	batch := make([]*api.Message, 0, len(p.playerStats))
	for id, stats := range p.playerStats {
		msg, err := p.prepareMessage(api.Code_SET_PLAYER_STATS, &api.SetPlayerStats{Id: id, Stats: stats})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
