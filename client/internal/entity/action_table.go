package entity

import (
	"container/ring"
	"fmt"
	"image/color"
	"time"
)

type ActionTable interface {
	AddConnecting(id uint64)
	AddDisconnecting(id uint64)
	AddKill(playerID uint64, playerClass Class, killerID uint64, killerClass Class)
	Get() []ActionMessage
}

type actionTable struct {
	messages *ring.Ring
	size     int
}

func newActionTable(n int) ActionTable {
	return &actionTable{
		messages: ring.New(n),
		size:     n,
	}
}

type ActionMessage struct {
	Message ActionMessageRender

	timestamp time.Time
	Alpha     uint8
}

type ActionMessageRender interface {
	Words() []WordRender
}

type WordRender interface {
	Text() string
	Color() color.NRGBA
}

type word struct {
	text  string
	color color.NRGBA
}

func (w word) Text() string       { return w.text }
func (w word) Color() color.NRGBA { return w.color }

func (t *actionTable) addMessage(msg ActionMessageRender) {
	t.messages.Value = ActionMessage{
		Message:   msg,
		timestamp: time.Now(),
	}
	t.messages = t.messages.Next()
}

type KillMessage struct {
	PlayerID    uint64
	IsPlayerBot bool
	PlayerClass Class

	KillerID    uint64
	IsKillerBot bool
	KillerClass Class

	cacheWords []WordRender
}

func (m *KillMessage) Words() []WordRender {
	if len(m.cacheWords) > 0 {
		return m.cacheWords
	}
	m.cacheWords = make([]WordRender, 0, 3)
	w := word{
		text:  fmt.Sprintf("%d", m.PlayerID),
		color: m.PlayerClass.Color(),
	}
	if m.IsPlayerBot {
		w.text = "bot " + w.text
	}
	m.cacheWords = append(m.cacheWords, w)
	m.cacheWords = append(m.cacheWords, word{
		text:  "kills",
		color: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF},
	})
	w = word{
		text:  fmt.Sprintf("%d", m.KillerID),
		color: m.KillerClass.Color(),
	}
	if m.IsKillerBot {
		w.text = "bot " + w.text
	}
	m.cacheWords = append(m.cacheWords, w)
	return m.cacheWords
}

func (t *actionTable) AddKill(playerID uint64, playerClass Class, killerID uint64, killerClass Class) {
	msg := &KillMessage{
		PlayerID:    playerID,
		PlayerClass: playerClass,
		KillerID:    killerID,
		KillerClass: killerClass,
	}
	if uint32(playerID>>32) == botPrefixID {
		msg.IsPlayerBot = true
		msg.PlayerID &= 0xFFFFFFFF
	}
	if uint32(killerID>>32) == botPrefixID {
		msg.IsKillerBot = true
		msg.KillerID &= 0xFFFFFFFF
	}

	t.addMessage(msg)
}

type ConnectingMessage struct {
	ID    uint64
	IsBot bool

	cacheWords []WordRender
}

func (m *ConnectingMessage) Words() []WordRender {
	if len(m.cacheWords) > 0 {
		return m.cacheWords
	}
	m.cacheWords = make([]WordRender, 0, 1)
	w := word{
		text:  fmt.Sprintf("%d connected", m.ID),
		color: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF},
	}
	if m.IsBot {
		w.text = "bot " + w.text
	}
	m.cacheWords = append(m.cacheWords, w)
	return m.cacheWords
}

func (t *actionTable) AddConnecting(id uint64) {
	msg := &ConnectingMessage{
		ID: id,
	}
	if uint32(id>>32) == botPrefixID {
		msg.IsBot = true
		msg.ID &= 0xFFFFFFFF
	}
	t.addMessage(msg)
}

type DisconnectingMessage struct {
	ID    uint64
	IsBot bool

	cacheWords []WordRender
}

func (m *DisconnectingMessage) Words() []WordRender {
	if len(m.cacheWords) > 0 {
		return m.cacheWords
	}
	m.cacheWords = make([]WordRender, 0, 1)
	w := word{
		text:  fmt.Sprintf("%d disconnected", m.ID),
		color: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF},
	}
	if m.IsBot {
		w.text = "bot " + w.text
	}
	m.cacheWords = append(m.cacheWords, w)
	return m.cacheWords
}

func (t *actionTable) AddDisconnecting(id uint64) {
	msg := &DisconnectingMessage{
		ID: id,
	}
	if uint32(id>>32) == botPrefixID {
		msg.IsBot = true
		msg.ID &= 0xFFFFFFFF
	}
	t.addMessage(msg)
}

const (
	durationActionMessage           float64 = 5.0 // seconds
	durationTransitionActionMessage float64 = 3.0 // seconds
)

func (t *actionTable) Get() []ActionMessage {
	out := make([]ActionMessage, 0, t.messages.Len())
	for iter, idx := t.messages.Prev(), 0; iter.Value != nil && idx < t.size; iter, idx = iter.Prev(), idx+1 {
		msg := iter.Value.(ActionMessage)
		if since := time.Since(msg.timestamp).Seconds(); since < durationTransitionActionMessage {
			msg.Alpha = 0
		} else {
			since -= durationTransitionActionMessage
			alpha := since / (durationActionMessage - durationTransitionActionMessage)
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
