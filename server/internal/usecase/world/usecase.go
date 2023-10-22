package usecase_world

import (
	"context"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/internal/log"
	"time"
)

type Usecase struct {
	logger     log.Logger
	ncProducer ncProducer
	repoWorld  repoWorld
	world      entity.World
}

type ncProducer interface {
	OnPlayerConnect(id uint64)
	OnPlayerDisconnect(id uint64)
	OnPlayerWasted(playerID uint64, playerClass entity.Class, killerID uint64, killerClass entity.Class)
	OnPlayerAttacked(id uint64, directionAim float64)
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

func NewUsecase(ctx context.Context, logger log.Logger, ncProducer ncProducer, repoWorld repoWorld) *Usecase {
	return &Usecase{
		logger:     logger.With(log.Module("uc_world")),
		ncProducer: ncProducer,
		repoWorld:  repoWorld,
		world: entity.NewWorld(ctx,
			time.Now().UnixNano(),
			entity.Point{X: entity.MaxWest, Y: entity.MaxNorth},
			entity.Point{X: entity.MaxEast, Y: entity.MaxSouth},
			ncProducer,
		),
	}
}
