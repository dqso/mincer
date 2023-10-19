package entity

import (
	"log"
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

func newBot(w World, id uint32, class Class) (Bot, func()) {
	b := &bot{
		Player: NewPlayer(uint64(botPrefixID)<<32|uint64(id), class),
		world:  w,
	}
	stop := make(chan struct{})

	go b.life(stop)

	return b, func() {
		close(stop)
	}
}

const botRespawnTime = time.Second * 5

func (b *bot) life(stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
		}
		//log.Printf("bot %d: New iteration my life", b.BotID())

		if b.IsDead() {
			log.Printf("bot %d: I'm dead. I'm waiting %s and I'm going to respawn", b.BotID(), botRespawnTime)
			time.Sleep(botRespawnTime)
			b.world.Respawn(b.GetPlayer())
			continue
		}

		// TODO
		if b.target == nil {
			//log.Printf("%d: need target", b.ID())
			b.target = b.world.SearchNearby(b.Position(), func(p Player) Player {
				if p.ID() == b.ID() {
					return nil
				}
				return p
			})
			if b.target != nil {
				log.Printf("%d: target has been found: %v", b.ID(), b.target)
			}
		}

		time.Sleep(time.Second)
	}
}

func (b *bot) GetPlayer() Player {
	return b.Player
}

func (b *bot) BotID() uint32 {
	return uint32(b.ID() & 0xFFFFFFFF)
}
