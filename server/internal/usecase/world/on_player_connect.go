package usecase_world

import (
	"github.com/dqso/mincer/server/internal/log"
	"log/slog"
)

func (uc *Usecase) OnPlayerConnect(connect chan uint64, disconnect chan uint64) {
	for {
		select {

		case id, ok := <-connect:
			if !ok {
				return
			}
			player, err := uc.world.NewPlayer(id)
			if err != nil {
				uc.logger.Error("unable to create the player",
					slog.Uint64("id", id),
					log.Err(err),
				)
				continue
			}
			uc.logger.Info("player has connected", slog.Uint64("id", id))
			uc.ncProducer.OnPlayerConnect(player.ID())
			uc.ncProducer.WorldInfo(id, uc.world)
			uc.ncProducer.PlayerList(id, uc.world.Players().Slice())
			uc.ncProducer.SpawnPlayer(player)

		case id, ok := <-disconnect:
			if !ok {
				return
			}
			uc.world.Players().Remove(id)
			uc.ncProducer.OnPlayerDisconnect(id)
			uc.logger.Info("player has disconnected", slog.Uint64("id", id))
		}
	}
}
