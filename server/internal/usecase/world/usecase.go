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
	PlayerList(toPlayerID uint64, players []entity.Player)
	SpawnPlayer(player entity.Player)
	SetPlayerClass(id uint64, class entity.Class)
	SetPlayerHP(id uint64, hp int64)
	SetPlayerRadius(id uint64, radius float64)
	SetPlayerSpeed(id uint64, speed float64)
	SetPlayerPosition(id uint64, position entity.Point)
}

func NewUsecase(ncProducer ncProducer) *Usecase {
	return &Usecase{
		ncProducer: ncProducer,
		world: entity.NewWorld(time.Now().UnixNano(),
			entity.Point{X: entity.MaxWest, Y: entity.MaxNorth},
			entity.Point{X: entity.MaxEast, Y: entity.MaxSouth},
		),
	}
}
