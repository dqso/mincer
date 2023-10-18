package usecase_world

import (
	"context"
	"log"
	"time"
)

func (uc Usecase) LifeCycle(ctx context.Context) chan struct{} {
	stopped := make(chan struct{})

	go func() {
		defer close(stopped)
		var lifeTime time.Duration
		const deltaTime = time.Millisecond * 10
		for {
			startPause := time.Now()
			select {
			case <-ctx.Done():
				return
			default:
			}

			changed := make(map[uint64]struct{})

			for _, player := range uc.world.Players().Slice() {
				if player.Move(deltaTime) {
					changed[player.ID()] = struct{}{}
				}
			}

			for id := range changed {
				player, ok := uc.world.Players().Get(id)
				if !ok {
					continue
				}
				uc.ncProducer.OnPlayerChange(player)
			}

			sleepTime := deltaTime - time.Since(startPause)
			if sleepTime < 0 {
				log.Print("WARN: sleep time is negative") // TODO
			}
			time.Sleep(sleepTime)
			lifeTime += time.Since(startPause)
		}
	}()

	return stopped
}
