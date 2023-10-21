package usecase_world

import (
	"context"
	"fmt"
)

func (uc *Usecase) ClientInfo(ctx context.Context, fromUserID uint64, direction float64, isMoving, attack bool, directionAim float64) error {
	p, ok := uc.world.Players().Get(fromUserID)
	if !ok {
		return fmt.Errorf("user not found") // TODO kick or everything else
	}
	p.SetDirection(direction, isMoving)
	p.SetAttack(attack, directionAim)
	return nil
}
