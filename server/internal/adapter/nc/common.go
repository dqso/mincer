package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
)

func (p *Producer) SendPayloadToClient(clientId uint64, payloadData []byte) error {
	// TODO log
	return p.server.SendPayloadToClient(clientId, payloadData)
}

func dtoPlayerToApiPlayer(player entity.Player) *api.Player {
	p := player.Position()
	return &api.Player{
		Id:     player.ID(),
		Class:  api.Class(player.Class()),
		Hp:     player.HP(),
		Radius: player.Radius(),
		Speed:  player.Speed(),
		X:      p.X,
		Y:      p.Y,
	}
}
