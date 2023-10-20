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
				player.Move(deltaTime, uc.world.SizeRect())
				player.Relax(deltaTime)
			}

			for _, player := range uc.world.Players().Slice() {
				if player.Attack() /* TODO && player.Class() == entity.ClassWarrior*/ { // TODO add cursor position for mage and ranger
					aPos := player.Position()
					rr := entity.DefaultAttackRadius * entity.DefaultAttackRadius
					uc.world.SearchNearby(player.Position(), func(p entity.Player) bool {
						if p.ID() == player.ID() {
							return false
						}
						pPos := p.Position()
						// ударять всех в радиусе
						if x, y := pPos.X-aPos.X, pPos.Y-aPos.Y; x*x+y*y <= rr {
							wasChanged := p.SetHP(p.HP() - int64(p.Power()))
							if p.HP() == 0 && wasChanged {
								uc.ncProducer.OnPlayerWasted(p.ID(), player.ID())
							}
						}
						return false
					})
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
