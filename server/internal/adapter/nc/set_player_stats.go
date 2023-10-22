package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/internal/log"
	"log/slog"
)

func (p *Producer) SetPlayerStats(id uint64, stats entity.PlayerStats) {
	p.mxPlayerStats.Lock()
	defer p.mxPlayerStats.Unlock()
	p.playerStats[id] = stats
}

func (p *Producer) setPlayerStatsBatch() []*api.Message {
	p.mxPlayerStats.Lock()
	defer p.mxPlayerStats.Unlock()
	defer clear(p.playerStats)
	batch := make([]*api.Message, 0, len(p.playerStats))
	for id, stats := range p.playerStats {
		p.logger.Debug("set player stats",
			slog.Uint64("id", id),
			log.Stats(stats),
		)
		batch = p.appendToBatch(batch, api.Code_SET_PLAYER_STATS, &api.SetPlayerStats{
			Id:    id,
			Stats: dtoPlayerStats(stats),
		})
	}
	return batch
}
