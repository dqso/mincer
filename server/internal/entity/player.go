package entity

import (
	"github.com/dqso/mincer/server/internal/api"
	"math"
	"sync"
	"time"
)

type Player interface {
	ID() uint64

	Class() Class
	SetClass(c Class)

	HP() int64
	SetHP(hp int64)

	Radius() float64
	SetRadius(r float64)

	Speed() float64
	SetSpeed(s float64)

	Position() Point
	SetPosition(p Point)

	SetDirection(direction float64, isMoving bool)
	Move(lifeTime time.Duration, mapSize Rect) (newPosition Point, wasMoved bool)
}

type player struct {
	id uint64

	mxStats sync.RWMutex
	class   Class
	hp      int64
	radius  float64
	speed   float64

	mxPosition sync.RWMutex
	x, y       float64
	// for calculating x and y
	direction float64
	isMoving  bool
}

func NewPlayer(id uint64, class Class) Player {
	return &player{
		id:     id,
		class:  class,
		hp:     defaultPlayerHP,
		radius: defaultPlayerRadius,
		speed:  defaultPlayerSpeed,
		x:      0,
		y:      0,
	}
}

func (p *player) ID() uint64 { return p.id }

func (p *player) Class() Class {
	p.mxStats.RLock()
	defer p.mxStats.RUnlock()
	return p.class
}

func (p *player) SetClass(c Class) {
	p.mxStats.Lock()
	defer p.mxStats.Unlock()
	p.class = c
}

func (p *player) HP() int64 {
	p.mxStats.RLock()
	defer p.mxStats.RUnlock()
	return p.hp
}

func (p *player) SetHP(hp int64) {
	p.mxStats.Lock()
	defer p.mxStats.Unlock()
	p.hp = hp
}

func (p *player) Radius() float64 {
	p.mxStats.RLock()
	defer p.mxStats.RUnlock()
	return p.radius
}

func (p *player) SetRadius(r float64) {
	p.mxStats.Lock()
	defer p.mxStats.Unlock()
	p.radius = r
}

func (p *player) Speed() float64 {
	p.mxStats.RLock()
	defer p.mxStats.RUnlock()
	return p.speed
}

func (p *player) SetSpeed(s float64) {
	p.mxStats.Lock()
	defer p.mxStats.Unlock()
	p.speed = s
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
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	if !p.isMoving {
		return Point{X: p.x, Y: p.y}, false
	}
	sin, cos := math.Sincos(p.direction * math.Pi / 180)
	x := p.x + p.speed*lifeTime.Seconds()*sin
	y := p.y - p.speed*lifeTime.Seconds()*cos
	if x < mapSize.LeftUp.X || mapSize.RightDown.X < x ||
		y < mapSize.LeftUp.Y || mapSize.RightDown.Y < y {
		return Point{X: p.x, Y: p.y}, false
	}
	p.x, p.y = x, y
	return Point{X: p.x, Y: p.y}, true
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
