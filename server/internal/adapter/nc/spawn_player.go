package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log/slog"
)

func (p *Producer) SpawnPlayer(player entity.Player) {
	p.mxSpawnPlayer.Lock()
	defer p.mxSpawnPlayer.Unlock()
	p.spawnPlayer[player.ID()] = player
}

func (p *Producer) spawnPlayerBatch() []*api.Message {
	p.mxSpawnPlayer.Lock()
	defer p.mxSpawnPlayer.Unlock()
	defer clear(p.spawnPlayer)
	batch := make([]*api.Message, 0, len(p.spawnPlayer))
	for id, player := range p.spawnPlayer {
		p.logger.Debug("spawn player",
			slog.Uint64("id", id),
			// TODO detailed info
		)
		batch = p.appendToBatch(batch, api.Code_SPAWN_PLAYER, &api.SpawnPlayer{
			Player: dtoPlayerToApiPlayer(player),
		})
	}
	return batch
}
