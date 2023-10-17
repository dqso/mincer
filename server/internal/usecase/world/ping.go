package usecase_world

import (
	"context"
)

func (uc Usecase) Ping(ctx context.Context, fromClientID uint64, ping string) error {
	if ping != "ping" {
		return nil
	}
	if err := uc.ncProducer.Pong(ctx, fromClientID, "pong"); err != nil {
		return err
	}
	return nil
}
