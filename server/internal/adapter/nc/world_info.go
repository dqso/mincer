package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) WorldInfo(toPlayerID uint64, world entity.World) {
	msg := &api.WorldInfo{
		Northwest: dtoPoint(world.Northwest()),
		Southeast: dtoPoint(world.Southeast()),
	}
	bts, err := p.marshalMessage(api.Code_WORLD_INFO, msg)
	if err != nil {
		log.Print(err) // TODO logger
		return
	}
	if err := p.SendPayloadToClient(toPlayerID, bts); err != nil {
		log.Print(err) // TODO logger
		return
	}
}
