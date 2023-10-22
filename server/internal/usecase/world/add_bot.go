package usecase_world

import (
	"github.com/dqso/mincer/server/internal/entity/ai"
	"log/slog"
)

func (uc *Usecase) AddBot() {
	class := uc.world.AcquireClass()
	bot, stopFunc := ai.NewBot(uc.logger, uc.world, uc.world.AcquireBotID(), class, uc.world.AcquireWeapon(class))
	uc.world.RegisterBot(bot, stopFunc)

	uc.ncProducer.OnPlayerConnect(bot.ID())
	uc.ncProducer.SpawnPlayer(bot.GetPlayer())
	uc.logger.Info("bot added",
		slog.Uint64("id", bot.ID()),
	)
}
