package entity

import (
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

type Me interface {
	Player
	GetPlayer() Player
	SetID(id uint64)

	Direction() (float64, bool)
	SetDirection(d float64, isMoving bool)
	Attack() bool
	SetAttack(v bool)
	CurrentCoolDown() float32
}

type me struct {
	Player

	direction float64
	isMoving  bool

	attack        bool
	isCoolDown    bool
	coolDownStart time.Time
}

func newEmptyMe() Me {
	return &me{
		Player:    newEmptyPlayer(),
		direction: 0,
	}
}

func (m *me) GetPlayer() Player { return m.Player }

func (m *me) SetID(id uint64)                       { m.Player.setID(id) }
func (m *me) Direction() (float64, bool)            { return m.direction, m.isMoving }
func (m *me) SetDirection(d float64, isMoving bool) { m.direction, m.isMoving = d, isMoving }
func (m *me) Attack() bool                          { return m.attack }

func (m *me) SetAttack(v bool) {
	m.attack = v
	if v && !m.isCoolDown {
		m.isCoolDown, m.coolDownStart = true, time.Now()
	}
}

func (m *me) CurrentCoolDown() float32 {
	if !m.isCoolDown {
		return 1000 // TODO from weapon m.MaxCoolDown()
	}
	current := float32(time.Since(m.coolDownStart).Seconds())
	if current >= 1000 /*TODO from weapon m.MaxCoolDown()*/ {
		m.isCoolDown = false
		return 1000 // TODO from weapon m.MaxCoolDown()
	}
	return current
}

func ColorBorderMe() color.Color {
	return colornames.White
}
