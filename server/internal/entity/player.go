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

	HP() int64
	IsDead() bool
	SetHP(hp int64)

	Position() Point
	SetPosition(p Point)

	SetDirection(direction float64, isMoving bool)
	Move(lifeTime time.Duration, mapSize Rect) (newPosition Point, wasMoved bool)
	SetAttack(a bool)
	Attack() bool
	Relax(lifeTime time.Duration)
}

type player struct {
	id uint64

	PlayerStats

	mxHP sync.RWMutex
	hp   int64

	mxPosition sync.RWMutex
	x, y       float64
	// for calculating x and y
	direction float64
	isMoving  bool

	mxAttack sync.Mutex
	isAttack bool
	coolDown float64
}

func NewPlayer(id uint64, class Class) Player {
	p := &player{
		id: id,
		PlayerStats: newPlayerStats(
			class,
			defaultPlayerRadius,
			defaultPlayerSpeed,
			defaultPlayerHP,
			defaultPlayerCoolDown,
			defaultPlayerPower,
		),
		x: 0,
		y: 0,
	}
	p.coolDown = p.MaxCoolDown()
	p.hp = p.MaxHP()
	return p
}

func (p *player) ID() uint64 { return p.id }

func (p *player) GetStats() PlayerStats {
	return p.PlayerStats
}

func (p *player) HP() int64 {
	p.mxHP.RLock()
	defer p.mxHP.RUnlock()
	return p.hp
}

func (p *player) IsDead() bool {
	return p.HP() <= 0
}

func (p *player) SetHP(hp int64) {
	if hp < 0 {
		hp = 0
	}
	if p.MaxHP() < hp {
		hp = 100
	}
	p.mxHP.Lock()
	defer p.mxHP.Unlock()
	p.hp = hp
}

func (p *player) Position() Point {
	p.mxPosition.RLock()
	defer p.mxPosition.RUnlock()
	return Point{X: p.x, Y: p.y}
}

func (p *player) SetPosition(point Point) {
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	p.x, p.y = point.X, point.Y
}

func (p *player) SetDirection(direction float64, isMoving bool) {
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	p.direction, p.isMoving = direction, isMoving
}

func (p *player) Move(lifeTime time.Duration, mapSize Rect) (newPosition Point, wasMoved bool) {
	if p.IsDead() {
		return Point{X: p.x, Y: p.y}, false
	}
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	if !p.isMoving {
		return Point{X: p.x, Y: p.y}, false
	}
	sin, cos := math.Sincos(p.direction * math.Pi / 180)
	x := p.x + p.Speed()*lifeTime.Seconds()*sin
	y := p.y - p.Speed()*lifeTime.Seconds()*cos
	if x < mapSize.LeftUp.X || mapSize.RightDown.X < x ||
		y < mapSize.LeftUp.Y || mapSize.RightDown.Y < y {
		return Point{X: p.x, Y: p.y}, false
	}
	p.x, p.y = x, y
	return Point{X: p.x, Y: p.y}, true
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
	if !p.isAttack || p.coolDown < p.MaxCoolDown() {
		return false
	}
	return true
}

func (p *player) Relax(lifeTime time.Duration) {
	p.mxAttack.Lock()
	defer p.mxAttack.Unlock()
	p.coolDown += p.MaxCoolDown() * lifeTime.Seconds()
}

type Class int32

const (
	ClassWarrior = Class(api.Class_WARRIOR)
	ClassMage    = Class(api.Class_MAGE)
	ClassRanger  = Class(api.Class_RANGER)
)

func Classes() []Class {
	return []Class{ClassWarrior, ClassMage, ClassRanger}
}
