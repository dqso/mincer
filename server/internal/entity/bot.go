package entity

import (
	"log"
	"math"
	"time"
)

type Bot interface {
	Player
	GetPlayer() Player
	BotID() uint32
}

const (
	botPrefixID = uint32(0xFFFFFFFF)
)

type bot struct {
	Player

	world  World
	target Player
}

func newBot(w World, id uint32, class Class, weapon Weapon) (Bot, func()) {
	b := &bot{
		Player: newPlayer(uint64(botPrefixID)<<32|uint64(id), class, weapon, w.Horn()),
		world:  w,
	}
	stop := make(chan struct{})

	go b.life(stop)

	return b, func() {
		close(stop)
	}
}

const (
	botRespawnTime     = time.Second * 5
	preferMageDistance = 100.0
)

func (b *bot) life(stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
		}
		//log.Printf("bot %d: New iteration my life", b.BotID())

		if b.IsDead() {
			//log.Printf("bot %d: I'm dead. I'm waiting %s and I'm going to respawn", b.BotID(), botRespawnTime)
			time.Sleep(botRespawnTime)
			b.world.Respawn(b.GetPlayer())
			b.target = nil
			continue
		}

		if b.target != nil {
			if b.target.IsDead() {
				b.target = nil
			}
		}

		// TODO
		if b.target == nil {
			//log.Printf("%d: need target", b.ID())
			b.target = b.world.SearchNearby(b.Position(), func(p Player) bool {
				if p.ID() == b.ID() {
					return false
				}
				return true
			})
			if b.target != nil {
				log.Printf("bot %d: target has been found: player %v %v", b.BotID(), b.target.ID(), b.target.Class()) // TODO logger
			}
		}

		if b.target != nil {
			var behavior botBehavior
			distance := b.Position().Distance(b.target.Position()) - b.target.Radius()
			if class := b.Class(); class == ClassWarrior {
				if distance > b.Weapon().AttackRadius()*0.97 {
					behavior = runToTarget
					// бежать до него
				} else {
					behavior = attackTarget
					// остановиться и лупить
				}
			} else if class == ClassMage || class == ClassRanger {
				if distance > DefaultMaxFireballDistance {
					behavior = runToTarget
					// бежать до него
				} else if distance < preferMageDistance {
					behavior = runAway
					// убегать от него
				} else {
					behavior = attackTarget
					// остановиться и лупить
				}
			}
			if behavior == runToTarget || behavior == runAway {
				botPos, vicPos := b.Position(), b.target.Position()
				direction := math.Mod(360.0-math.Atan2(botPos.X-vicPos.X, botPos.Y-vicPos.Y)*180.0/math.Pi, 360.0)
				if behavior == runAway {
					direction = math.Mod(direction+180.0, 360.0)
				}
				b.SetDirection(direction, true)
				b.SetAttack(false)
			} else if behavior == attackTarget {
				b.SetDirection(0.0, false)
				b.SetAttack(true)
			}
		}

		time.Sleep(time.Millisecond * 300)
	}
}

type botBehavior uint8

const (
	idle botBehavior = iota
	runAway
	runToTarget
	attackTarget
)

func (b *bot) GetPlayer() Player {
	return b.Player
}

func (b *bot) BotID() uint32 {
	return uint32(b.ID() & 0xFFFFFFFF)
}
