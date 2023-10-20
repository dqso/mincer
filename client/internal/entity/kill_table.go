package entity

import (
	"container/ring"
	"time"
)

type KillTable interface {
	Add(playerID uint64, player Player, killerID uint64, killer Player)
	Get() []KillMessage
}

type killTable struct {
	messages *ring.Ring
	size     int
}

func newKillTable(n int) KillTable {
	return &killTable{
		messages: ring.New(n),
		size:     n,
	}
}

type KillMessage struct {
	PlayerID    uint64
	IsPlayerBot bool
	PlayerClass Class

	KillerID    uint64
	IsKillerBot bool
	KillerClass Class

	timestamp time.Time
	Alpha     uint8
}

func (t *killTable) Add(playerID uint64, player Player, killerID uint64, killer Player) {
	msg := KillMessage{
		PlayerID: playerID,
		KillerID: killerID,

		timestamp: time.Now(),
	}
	if uint32(playerID>>32) == botPrefixID {
		msg.IsPlayerBot = true
		msg.PlayerID &= 0xFFFFFFFF
	}
	if uint32(killerID>>32) == botPrefixID {
		msg.IsKillerBot = true
		msg.KillerID &= 0xFFFFFFFF
	}
	if player != nil {
		msg.PlayerClass = player.Class()
	}
	if killer != nil {
		msg.KillerClass = killer.Class()
	}

	t.messages.Value = msg
	t.messages = t.messages.Next()
}

const (
	durationKillMessage           float64 = 5.0 // seconds
	durationTransitionKillMessage float64 = 3.0 // seconds
)

func (t *killTable) Get() []KillMessage {
	out := make([]KillMessage, 0, t.messages.Len())
	for iter, idx := t.messages.Prev(), 0; iter.Value != nil && idx < t.size; iter, idx = iter.Prev(), idx+1 {
		msg := iter.Value.(KillMessage)
		if since := time.Since(msg.timestamp).Seconds(); since < durationTransitionKillMessage {
			msg.Alpha = 0
		} else {
			since -= durationTransitionKillMessage
			alpha := since / (durationKillMessage - durationTransitionKillMessage)
			if alpha >= 1 {
				break
			} else {
				msg.Alpha = uint8(alpha * 255)
			}
		}
		out = append(out, msg)
	}
	return out
}
