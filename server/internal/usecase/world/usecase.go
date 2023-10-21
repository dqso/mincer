package usecase_world

import (
	"github.com/dqso/mincer/server/internal/entity"
	"time"
)

type Usecase struct {
	ncProducer ncProducer
	repoWorld  repoWorld
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
	SetPlayerWeapon(id uint64, w entity.Weapon)
	CreateProjectile(projectile entity.Projectile)
	SetProjectilePosition(id uint64, position entity.Point)
	DeleteProjectile(id uint64)
}

type repoWorld interface {
	AcquireProjectileID() uint64
}

func NewUsecase(ncProducer ncProducer, repoWorld repoWorld) *Usecase {
	return &Usecase{
		ncProducer: ncProducer,
		repoWorld:  repoWorld,
		world: entity.NewWorld(time.Now().UnixNano(),
			entity.Point{X: entity.MaxWest, Y: entity.MaxNorth},
			entity.Point{X: entity.MaxEast, Y: entity.MaxSouth},
			ncProducer,
		),
	}
}
