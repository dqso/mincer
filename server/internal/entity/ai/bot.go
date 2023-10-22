package ai

import (
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/internal/log"
	"log/slog"
	"math"
	"time"
)

const (
	botPrefixID        = uint32(0xFFFFFFFF)
	botRespawnTime     = time.Second * 5
	preferMageDistance = 100.0
)

type Bot struct {
	logger log.Logger
	entity.Player

	world  entity.World
	target entity.Player
}

func NewBot(logger log.Logger, w entity.World, id uint32, class entity.Class, weapon entity.Weapon) (*Bot, func()) {
	b := &Bot{
		logger: logger,
		Player: entity.NewPlayer(uint64(botPrefixID)<<32|uint64(id), class, weapon, w.Horn()),
		world:  w,
	}
	stop := make(chan struct{})

	go b.life(stop)

	return b, func() {
		close(stop)
	}
}

func (b *Bot) GetPlayer() entity.Player { return b.Player }
func (b *Bot) BotID() uint32            { return uint32(b.ID() & 0xFFFFFFFF) }

// TODO move to usecase lifeCycle
func (b *Bot) life(stop chan struct{}) {
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
			b.target = b.world.SearchNearby(b.Position(), func(p entity.Player) bool {
				if p.ID() == b.ID() || p.IsDead() {
					return false
				}
				return true
			})
			if b.target != nil {
				b.logger.Debug("target has been found: player %v %v",
					slog.Uint64("bot_id", b.ID()),
					slog.Uint64("target", b.target.ID()),
				)
			}
		}

		if b.target != nil {
			var behavior botBehavior
			distance := b.Position().Distance(b.target.Position()) - b.target.Radius()
			// анализ поведений
			if class := b.Class(); class == entity.ClassWarrior {
				if distance > b.Weapon().AttackRadius()*0.97 {
					behavior = runToTarget | noAttack
					// бежать до него
				} else {
					behavior = stay | attackTarget
					// остановиться и лупить
				}
			} else if class == entity.ClassMage || class == entity.ClassRanger {
				if distance > entity.DefaultMaxFireballDistance {
					behavior = runToTarget | noAttack
					// бежать до него
				} else if distance < preferMageDistance {
					behavior = runAway | attackTarget
					// убегать от него и стрелять
				} else {
					behavior = stay | attackTarget
					// остановиться и лупить
				}
			}

			// выполнение поведений
			botPos, vicPos := b.Position(), b.target.Position()
			direction := math.Mod(360.0-math.Atan2(botPos.X-vicPos.X, botPos.Y-vicPos.Y)*180.0/math.Pi, 360.0) // TODO create a new method in the Point struct
			if behavior.Or(runToTarget, runAway) {
				direction := direction // сохранить внешнее направление для последующего поведения
				if behavior.Or(runAway) {
					// развернуться
					direction = math.Mod(direction+180.0, 360.0)
				}
				b.SetDirection(direction, true)
			} else if behavior.Or(stay) {
				b.SetDirection(0.0, false)
			}
			if behavior.Or(attackTarget) {
				b.SetAttack(true, direction)
			} else if behavior.Or(noAttack) {
				b.SetAttack(false, direction)
			}
		}

		time.Sleep(time.Millisecond * 300)
	}
}

type botBehavior uint8

func (b botBehavior) Or(behaviors ...botBehavior) bool {
	for _, behavior := range behaviors {
		if b&behavior > 0 {
			return true
		}
	}
	return false
}

const (
	_ botBehavior = 1 << iota
	stay
	runAway
	runToTarget
	noAttack
	attackTarget
)
