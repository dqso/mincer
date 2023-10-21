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
				middlePosition := oldPosition.Middle(newPosition)
				victim := projectile.CollisionAnalysis(middlePosition, players)
				if victim != nil {
					if projectile.AttackRadius() <= 1e-3 {
						victim.DealDamage(projectile.Owner(), projectile.Damage())
					} else {
						uc.dealDamageInRadius(
							middlePosition, projectile.Owner(),
							projectile.AttackRadius(), projectile.Damage(),
						)
					}
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
						weapon := player.Weapon()
						uc.dealDamageInRadius(
							player.Position(), player.ID(),
							weapon.AttackRadius(), weapon.Damage(),
						)
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

func (uc *Usecase) dealDamageInRadius(position entity.Point, attacker uint64, radius float64, damage entity.Damage) {
	uc.world.SearchNearby(position, func(p entity.Player) bool {
		if p.ID() == attacker {
			return false
		}
		pPos := p.Position()
		rr := radius + p.Radius() // учитывать радиус врага, а не только его центр тела
		rr = rr * rr
		// ударять всех в радиусе
		if x, y := pPos.X-position.X, pPos.Y-position.Y; x*x+y*y <= rr {
			p.DealDamage(attacker, damage)
		}
		return false
	})
}
