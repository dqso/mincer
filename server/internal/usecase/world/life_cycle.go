package usecase_world

import (
	"context"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
	"time"
)

func (uc *Usecase) LifeCycle(ctx context.Context) chan struct{} {
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
				if newPos, wasMoved := player.Move(deltaTime, uc.world.SizeRect()); wasMoved {
					uc.ncProducer.SetPlayerPosition(player.ID(), newPos)
				}
				player.Relax(deltaTime)
			}

			for _, player := range uc.world.Players().Slice() {
				if player.Attack() { // TODO add cursor position for mage and ranger
					log.Print(player.ID(), "attack")
					p := uc.world.SearchNearby(player.Position(), func(p entity.Player) entity.Player {
						if p.ID() == player.ID() {
							return nil
						}
						return p
					}) // TODO ударять всех
					if p != nil {
						p.SetHP(p.HP() - int64(p.Power()))
						uc.ncProducer.SetPlayerHP(p.ID(), p.HP())
					}
				}
			}

			for _, player := range uc.world.Players().Slice() {
				if player.HP() <= 0 {
					player.SetHP(0)
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
