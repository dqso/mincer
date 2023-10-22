package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
)

func (p *Producer) WorldInfo(toPlayerID uint64, world entity.World) {
	p.SendPayloadToClient(toPlayerID, api.Code_WORLD_INFO, &api.WorldInfo{
		Northwest: dtoPoint(world.Northwest()),
		Southeast: dtoPoint(world.Southeast()),
	})
}
