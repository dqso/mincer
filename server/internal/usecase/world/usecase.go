package usecase_world

import (
	"context"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
	"time"
)

type Usecase struct {
	ncProducer ncProducer
	world      entity.World
}

type ncProducer interface {
	Pong(ctx context.Context, toClientID uint64, pong string) error
	OnPlayerConnect(id uint64)
	OnPlayerDisconnect(id uint64)
	DirectPlayerList(toPlayerID uint64, players []entity.Player)
	OnPlayerChange(player entity.Player)
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

func (uc *Usecase) OnPlayerConnect(connect chan uint64, disconnect chan uint64) {
	for {
		select {

		case id, ok := <-connect:
			if !ok {
				return
			}
			player, err := uc.world.AddPlayer(id)
			if err != nil {
				log.Print(err) // TODO logger
				continue
			}
			uc.ncProducer.OnPlayerConnect(player.ID())
			uc.ncProducer.DirectPlayerList(id, uc.world.Players().Slice())
			uc.ncProducer.OnPlayerChange(player)

		case id, ok := <-disconnect:
			if !ok {
				return
			}
			uc.world.RemovePlayer(id)
			uc.ncProducer.OnPlayerDisconnect(id)
		}
	}
}
