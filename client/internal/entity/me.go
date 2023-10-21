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
	Attack() (bool, float64)
	SetAttack(v bool, direction float64)
	CurrentCoolDown() float64
}

type me struct {
	Player

	direction float64
	isMoving  bool

	attack        bool
	directionAim  float64
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
func (m *me) Attack() (bool, float64)               { return m.attack, m.directionAim }

func (m *me) SetAttack(v bool, direction float64) {
	if m.HP() == 0 {
		return
	}
	m.attack, m.directionAim = v, direction
	if v && !m.isCoolDown {
		m.isCoolDown, m.coolDownStart = true, time.Now()
	}
}

func (m *me) CurrentCoolDown() float64 {
	maxCoolDown := m.Weapon().CoolDown()
	if !m.isCoolDown {
		return maxCoolDown
	}
	current := time.Since(m.coolDownStart).Seconds()
	if current >= maxCoolDown {
		m.isCoolDown = false
		return maxCoolDown
	}
	return current
}

func ColorBorderMe() color.Color {
	return colornames.White
}
