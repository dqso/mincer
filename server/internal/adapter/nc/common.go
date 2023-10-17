package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
)

func dtoPlayerToPublicPlayer(player entity.Player) *api.PublicPlayer {
	x, y := player.Position()
	hp, radius, dead := player.PublicStats()
	return &api.PublicPlayer{
		PlayerId: player.ID(),
		X:        x,
		Y:        y,
		Hp:       hp,
		Radius:   radius,
		Dead:     dead,
	}
}
