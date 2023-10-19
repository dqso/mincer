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

			for _, player := range uc.world.Players().Slice() {
				if newPos, wasMoved := player.Move(deltaTime); wasMoved {
					uc.ncProducer.SetPlayerPosition(player.ID(), newPos)
				}
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
