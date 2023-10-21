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

			players := uc.world.Players().Slice()

			for _, projectile := range uc.world.ProjectileList().Slice() {
				//log.Printf("projectile %d: %v", projectile.ID(), projectile.Position())
				oldPosition := projectile.Position()
				newPosition, outOfRange := projectile.Move(deltaTime, uc.world.SizeRect())
				victim := projectile.CollisionAnalysis(oldPosition, newPosition, players)
				if victim != nil {
					victim.SetHP(victim.HP() - projectile.PhysicalDamage() - projectile.MagicalDamage())
					// TODO взрыв фаербола на всех близлежащих игроков в радиусе
				}
				if outOfRange || victim != nil {
					uc.world.ProjectileList().Remove(projectile.ID())
					uc.ncProducer.DeleteProjectile(projectile.ID())
				}
			}

			for _, player := range players {
				player.Move(deltaTime, uc.world.SizeRect())
				player.Relax(deltaTime)
			}

			for _, player := range players {
				if player.Attack() {
					if projectile, isMelee := player.Weapon().Attack(player, player.DirectionAim(), uc.repoWorld); isMelee {
						aPos := player.Position()
						uc.world.SearchNearby(player.Position(), func(p entity.Player) bool {
							if p.ID() == player.ID() {
								return false
							}
							pPos := p.Position()
							rr := player.Weapon().AttackRadius() + p.Radius() // учитывать радиус врага, а не только его центр тела
							rr = rr * rr
							// ударять всех в радиусе
							if x, y := pPos.X-aPos.X, pPos.Y-aPos.Y; x*x+y*y <= rr {
								wasChanged := p.SetHP(p.HP() - p.Weapon().PhysicalDamage() - p.Weapon().MagicalDamage())
								if p.HP() == 0 && wasChanged {
									uc.ncProducer.OnPlayerWasted(p.ID(), player.ID())
								}
							}
							return false
						})
					} else if projectile != nil {
						uc.world.ProjectileList().Add(projectile)
						uc.ncProducer.CreateProjectile(projectile)
						//log.Printf("created projectile %d: %+v", projectile.ID(), projectile)
					}
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
