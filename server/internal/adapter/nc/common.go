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
		Id:    player.ID(),
		Stats: dtoPlayerStats(player.GetStats()),
		Hp:    player.HP(),
		X:     p.X,
		Y:     p.Y,
	}
}

func dtoPlayerStats(stats entity.PlayerStats) *api.PlayerStats {
	return &api.PlayerStats{
		Class:       api.Class(stats.Class()),
		Radius:      stats.Radius(),
		Speed:       stats.Speed(),
		MaxHP:       stats.MaxHP(),
		MaxCoolDown: stats.MaxCoolDown(),
		Power:       stats.Power(),
	}
}
