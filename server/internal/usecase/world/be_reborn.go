package usecase_world

import (
	"context"
	"fmt"
	"log/slog"
)

func (uc *Usecase) BeReborn(ctx context.Context, fromUserID uint64) error {
	player, ok := uc.world.Players().Get(fromUserID)
	if !ok {
		return fmt.Errorf("player not found")
	}
	if player.HP() != 0 {
		return fmt.Errorf("player has some hp")
	}
	uc.world.Respawn(player)
	uc.logger.Info("player has been reborn",
		slog.Uint64("id", fromUserID),
	)
	return nil
}
