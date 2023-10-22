package usecase_world

import (
	"context"
	"log/slog"
)

func (uc *Usecase) Quit(ctx context.Context, fromUserID uint64) error {
	uc.world.Players().Remove(fromUserID)
	uc.ncProducer.OnPlayerDisconnect(fromUserID)
	uc.logger.Info("player has left",
		slog.Uint64("id", fromUserID),
	)
	return nil
}
