package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/internal/log"
	"google.golang.org/protobuf/proto"
	"image/color"
	"log/slog"
)

func (p *Producer) SendPayloadToClient(clientId uint64, code api.Code, msg proto.Message) bool {
	bts, err := p.marshalMessage(code, msg)
	if err != nil {
		p.logger.Error("unable to marshal the message for the client",
			slog.String("code", code.String()),
			slog.Uint64("client_id", clientId),
			log.Err(err),
		)
		return false
	}
	err = p.server.SendPayloadToClient(clientId, bts)
	if err != nil {
		p.logger.Error("unable to send the message to client",
			slog.String("code", code.String()),
			slog.Uint64("client_id", clientId),
			slog.Int("size_message", len(bts)),
			log.Err(err),
		)
		return false
	}
	p.logger.Debug("the message has been sent to the client",
		slog.String("code", code.String()),
		slog.Uint64("client_id", clientId),
	)
	return true
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

func dtoResist(r entity.Resist) *api.Resist {
	return &api.Resist{
		Physical: r.PhysicalResist(),
		Magical:  r.MagicalResist(),
	}
}

func dtoPlayerStats(stats entity.PlayerStats) *api.PlayerStats {
	return &api.PlayerStats{
		Class:  api.Class(stats.Class()),
		Resist: dtoResist(stats),
		Radius: stats.Radius(),
		Speed:  stats.Speed(),
		MaxHP:  stats.MaxHP(),
	}
}

func dtoDamage(d entity.Damage) *api.Damage {
	return &api.Damage{
		Physical: d.Physical(),
		Magical:  d.Magical(),
	}
}

func dtoWeapon(weapon entity.Weapon) *api.Weapon {
	return &api.Weapon{
		Name:     weapon.Name(),
		Damage:   dtoDamage(weapon.Damage()),
		CoolDown: weapon.CoolDown(),
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
