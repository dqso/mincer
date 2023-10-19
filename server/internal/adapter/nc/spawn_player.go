package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
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
	for _, player := range p.spawnPlayer {
		msg, err := p.prepareMessage(api.Code_SPAWN_PLAYER, &api.SpawnPlayer{Player: dtoPlayerToApiPlayer(player)})
		if err != nil {
			log.Print(err) // TODO logger
			continue
		}
		batch = append(batch, msg)
	}
	return batch
}
