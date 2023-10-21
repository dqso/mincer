package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"image/color"
)

func (p *Producer) SendPayloadToClient(clientId uint64, payloadData []byte) error {
	// TODO log
	return p.server.SendPayloadToClient(clientId, payloadData)
}

func dtoPlayerToApiPlayer(player entity.Player) *api.Player {
	p := player.Position()
	return &api.Player{
		Id:     player.ID(),
		Stats:  dtoPlayerStats(player.GetStats()),
		Weapon: dtoWeapon(player.Weapon()),
		Hp:     player.HP(),
		X:      p.X,
		Y:      p.Y,
	}
}

func dtoPlayerStats(stats entity.PlayerStats) *api.PlayerStats {
	return &api.PlayerStats{
		Class:  api.Class(stats.Class()),
		Radius: stats.Radius(),
		Speed:  stats.Speed(),
		MaxHP:  stats.MaxHP(),
	}
}

func dtoWeapon(weapon entity.Weapon) *api.Weapon {
	return &api.Weapon{
		Name:           weapon.Name(),
		PhysicalDamage: weapon.Damage().Physical(),
		MagicalDamage:  weapon.Damage().Magical(),
		CoolDown:       weapon.CoolDown(),
	}
}

func dtoProjectile(p entity.Projectile) *api.Projectile {
	return &api.Projectile{
		Id:        p.ID(),
		Color:     dtoColor(p.Color()),
		Position:  dtoPoint(p.Position()),
		Radius:    p.Radius(),
		Speed:     p.Speed(),
		Direction: p.Direction(),
	}
}

func dtoPoint(p entity.Point) *api.Point {
	return &api.Point{
		X: p.X,
		Y: p.Y,
	}
}

func dtoColor(p color.NRGBA) *api.Color {
	return &api.Color{
		Rgba: uint32(p.R)<<24 | uint32(p.G)<<16 | uint32(p.B)<<8 | uint32(p.A),
	}
}
