package entity

import (
	"math"
	"sync"
	"time"
)

type Player interface {
	ID() uint64
	SetHP(hp int64)
	SetPosition(x, y float64)
	Position() Point
	SetDirection(direction float64, isMoving bool)
	Move(lifeTime time.Duration) (isMoved bool)
	Radius() float64
	PublicStats() (hp int64, radius float64, dead bool)
}

type player struct {
	clientID uint64

	x, y       float64
	direction  float64
	isMoving   bool
	mxPosition sync.RWMutex

	hp      int64
	radius  float64
	speed   float64
	dead    bool
	mxStats sync.RWMutex
}

func NewPlayer(id uint64) Player {
	p := &player{
		clientID: id,
		x:        0,
		y:        0,
	}
	p.SetHP(defaultPlayerHP)
	return p
}

func (p *player) SetHP(hp int64) {
	p.mxStats.Lock()
	defer p.mxStats.Unlock()
	p.hp = hp
	if hp >= 100 {
		p.dead = false
		p.radius, p.speed = minPlayerRadius, maxPlayerSpeed
	} else if hp <= 0 {
		p.dead = true
		p.radius, p.speed = maxPlayerRadius, minPlayerSpeed
	} else {
		r := 1 / float64(hp)
		p.dead = false
		p.radius = (maxPlayerRadius-minPlayerRadius)*r + minPlayerRadius
		p.speed = maxPlayerSpeed - (maxPlayerSpeed-minPlayerSpeed)*r
	}
}

func (p *player) ID() uint64 {
	return p.clientID
}

func (p *player) SetPosition(x float64, y float64) {
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	p.x, p.y = x, y
}

func (p *player) Position() Point {
	p.mxPosition.RLock()
	defer p.mxPosition.RUnlock()
	return Point{X: p.x, Y: p.y}
}

func (p *player) SetDirection(direction float64, isMoving bool) {
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	p.direction, p.isMoving = direction, isMoving
}

func (p *player) Move(lifeTime time.Duration) (isMoved bool) {
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	if !p.isMoving {
		return false
	}
	sin, cos := math.Sincos(p.direction * math.Pi / 180)
	p.x = p.x + p.speed*lifeTime.Seconds()*sin
	p.y = p.y - p.speed*lifeTime.Seconds()*cos
	return true
}

func (p *player) Radius() float64 {
	p.mxStats.RLock()
	defer p.mxStats.RUnlock()
	return p.radius
}

func (p *player) PublicStats() (hp int64, radius float64, dead bool) {
	p.mxStats.RLock()
	defer p.mxStats.RUnlock()
	return p.hp, p.radius, p.dead
}

type Players interface {
	Get(id uint64) (Player, bool)
	Add(p Player)
	Remove(id uint64)
	Slice() []Player
}

type players struct {
	byID   map[uint64]Player
	mxByID sync.RWMutex
}

func NewPlayers() Players {
	return &players{
		byID: make(map[uint64]Player),
	}
}

func (pp *players) Get(id uint64) (Player, bool) {
	pp.mxByID.RLock()
	defer pp.mxByID.RUnlock()
	p, ok := pp.byID[id]
	return p, ok
}

func (pp *players) Add(p Player) {
	pp.mxByID.Lock()
	defer pp.mxByID.Unlock()
	pp.byID[p.ID()] = p
}

func (pp *players) Remove(id uint64) {
	pp.mxByID.Lock()
	defer pp.mxByID.Unlock()
	delete(pp.byID, id)
}

func (pp *players) Slice() []Player {
	pp.mxByID.RLock()
	defer pp.mxByID.RUnlock()
	out := make([]Player, 0, len(pp.byID))
	for _, p := range pp.byID {
		out = append(out, p)
	}
	return out
}
