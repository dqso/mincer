package usecase_world

import (
	"github.com/dqso/mincer/server/internal/entity"
	"time"
)

type Usecase struct {
	ncProducer ncProducer
	world      entity.World
}

type ncProducer interface {
	OnPlayerConnect(id uint64)
	OnPlayerDisconnect(id uint64)
	OnPlayerWasted(id uint64, killer uint64)
	WorldInfo(toPlayerID uint64, world entity.World)
	PlayerList(toPlayerID uint64, players []entity.Player)
	SpawnPlayer(player entity.Player)
	SetPlayerStats(id uint64, stats entity.PlayerStats)
	SetPlayerHP(id uint64, hp int32)
	SetPlayerPosition(id uint64, position entity.Point)
}

func NewUsecase(ncProducer ncProducer) *Usecase {
	return &Usecase{
		ncProducer: ncProducer,
		world: entity.NewWorld(time.Now().UnixNano(),
			entity.Point{X: entity.MaxWest, Y: entity.MaxNorth},
			entity.Point{X: entity.MaxEast, Y: entity.MaxSouth},
			ncProducer,
		),
	}
}
