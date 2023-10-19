package usecase_world

import (
	"context"
)

func (uc *Usecase) Quit(ctx context.Context, fromUserID uint64) error {
	uc.world.Players().Remove(fromUserID)
	return nil
}
