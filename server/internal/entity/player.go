package entity

import (
	"github.com/dqso/mincer/server/internal/api"
	"math"
	"sync"
	"time"
)

type Player interface {
	ID() uint64

	GetStats() PlayerStats
	PlayerStats

	Weapon() Weapon

	HP() int32
	IsDead() bool
	SetHP(hp int32) (wasChanged bool)

	Position() Point
	SetPosition(p Point)

	SetDirection(direction float64, isMoving bool)
	Move(lifeTime time.Duration, mapSize Rect) (wasMoved bool)
	SetAttack(a bool)
	Attack() bool
	Relax(lifeTime time.Duration)
}

type player struct {
	horn Horn

	id uint64

	PlayerStats

	weapon Weapon

	mxHP sync.RWMutex
	hp   int32

	mxPosition sync.RWMutex
	x, y       float64
	// for calculating x and y
	direction float64
	isMoving  bool

	mxAttack sync.Mutex
	isAttack bool
	coolDown float64
}

func newPlayer(id uint64, class Class, weapon Weapon, horn Horn) Player {
	p := &player{
		horn:   horn,
		id:     id,
		weapon: weapon,
		PlayerStats: newPlayerStats(
			class,
			defaultPlayerRadius,
			defaultPlayerSpeed,
			defaultPlayerHP,
		),
		x: 0,
		y: 0,
	}
	p.coolDown = weapon.CoolDown()
	p.hp = p.MaxHP()
	return p
}

func (p *player) ID() uint64 { return p.id }

func (p *player) GetStats() PlayerStats {
	return p.PlayerStats
}

func (p *player) Weapon() Weapon {
	return p.weapon
}

func (p *player) SetClass(v Class) {
	p.PlayerStats.SetClass(v)
	p.horn.SetPlayerStats(p.ID(), p.PlayerStats)
}

func (p *player) SetRadius(v float64) {
	p.PlayerStats.SetRadius(v)
	p.horn.SetPlayerStats(p.ID(), p.PlayerStats)
}

func (p *player) SetSpeed(v float64) {
	p.PlayerStats.SetSpeed(v)
	p.horn.SetPlayerStats(p.ID(), p.PlayerStats)
}

func (p *player) SetMaxHP(v int32) {
	p.PlayerStats.SetMaxHP(v)
	p.horn.SetPlayerStats(p.ID(), p.PlayerStats)
}

func (p *player) HP() int32 {
	p.mxHP.RLock()
	defer p.mxHP.RUnlock()
	return p.hp
}

func (p *player) IsDead() bool {
	return p.HP() <= 0
}

func (p *player) SetHP(hp int32) (wasChanged bool) {
	if hp < 0 {
		hp = 0
	}
	if p.MaxHP() < hp {
		hp = 100
	}
	func() {
		p.mxHP.Lock()
		defer p.mxHP.Unlock()
		if p.hp != hp {
			wasChanged = true
			p.hp = hp
		}
	}()
	p.horn.SetPlayerHP(p.ID(), hp)
	return
}

func (p *player) Position() Point {
	p.mxPosition.RLock()
	defer p.mxPosition.RUnlock()
	return Point{X: p.x, Y: p.y}
}

func (p *player) SetPosition(point Point) {
	func() {
		p.mxPosition.Lock()
		defer p.mxPosition.Unlock()
		p.x, p.y = point.X, point.Y
	}()
	p.horn.SetPlayerPosition(p.ID(), Point{p.x, p.y})
}

func (p *player) SetDirection(direction float64, isMoving bool) {
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	p.direction, p.isMoving = direction, isMoving
}

func (p *player) Move(lifeTime time.Duration, mapSize Rect) (wasMoved bool) {
	if p.IsDead() {
		wasMoved = false
		return
	}
	func() {
		p.mxPosition.Lock()
		defer p.mxPosition.Unlock()
		if !p.isMoving {
			wasMoved = false
			return
		}
		sin, cos := math.Sincos(p.direction * math.Pi / 180)
		x := p.x + p.Speed()*lifeTime.Seconds()*sin
		y := p.y - p.Speed()*lifeTime.Seconds()*cos
		if x < mapSize.LeftUp.X || mapSize.RightDown.X < x ||
			y < mapSize.LeftUp.Y || mapSize.RightDown.Y < y {
			wasMoved = false
			return
		}
		p.x, p.y = x, y
		wasMoved = true
		return
	}()
	if wasMoved {
		p.horn.SetPlayerPosition(p.ID(), Point{X: p.x, Y: p.y})
	}
	return wasMoved
}

func (p *player) SetAttack(a bool) {
	p.mxAttack.Lock()
	defer p.mxAttack.Unlock()
	p.isAttack = a
}

func (p *player) Attack() (isAllowed bool) {
	if p.HP() <= 0 {
		return false
	}
	p.mxAttack.Lock()
	defer p.mxAttack.Unlock()
	defer func() {
		if isAllowed {
			p.isAttack = false
			p.coolDown = 0
		}
	}()
	if !p.isAttack || p.coolDown < p.weapon.CoolDown() {
		return false
	}
	return true
}

func (p *player) Relax(lifeTime time.Duration) {
	p.mxAttack.Lock()
	defer p.mxAttack.Unlock()
	p.coolDown += lifeTime.Seconds()
}

type Class int32

const (
	ClassWarrior = Class(api.Class_WARRIOR)
	ClassMage    = Class(api.Class_MAGE)
	ClassRanger  = Class(api.Class_RANGER)
)

func Classes() []Class {
	return []Class{
		ClassWarrior,
		//ClassMage,
		//ClassRanger,
	}
}
