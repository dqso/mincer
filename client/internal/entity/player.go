package entity

import (
	"sort"
	"sync"
)

type Me interface {
	Player
	GetPlayer() Player
	SetID(id uint64)

	Direction() float64
	SetDirection(d float64)
	Speed() float64
}

type me struct {
	Player

	direction float64
	speed     float64
}

func newEmptyMe() Me {
	return &me{
		Player:    newEmptyPlayer(),
		direction: 0,
		speed:     0,
	}
}

func (m *me) GetPlayer() Player      { return m.Player }
func (m *me) SetID(id uint64)        { m.Player.setID(id) }
func (m *me) Direction() float64     { return m.direction }
func (m *me) SetDirection(d float64) { m.direction = d }
func (m *me) Speed() float64         { return m.speed }

type Player interface {
	ID() uint64
	setID(id uint64)
	Position() (float64, float64)
	SetPosition(x, y float64)
	PositionFloat32() (float32, float32)
	SetStats(hp int64, radius float64, dead bool)
	Radius() float64
	RadiusFloat32() float32
}

type player struct {
	id uint64

	x, y float64

	hp     int64
	radius float64
	dead   bool
}

func newEmptyPlayer() Player {
	return &player{}
}

func NewPlayer(id uint64, x, y float64, hp int64, radius float64, dead bool) Player {
	return &player{
		id:     id,
		x:      x,
		y:      y,
		hp:     hp,
		radius: radius,
		dead:   dead,
	}
}

func (p *player) ID() uint64                          { return p.id }
func (p *player) setID(id uint64)                     { p.id = id }
func (p *player) Position() (float64, float64)        { return p.x, p.y }
func (p *player) SetPosition(x, y float64)            { p.x, p.y = x, y }
func (p *player) PositionFloat32() (float32, float32) { return float32(p.x), float32(p.y) }
func (p *player) Radius() float64                     { return p.radius }
func (p *player) RadiusFloat32() float32              { return float32(p.radius) }
func (p *player) SetStats(hp int64, radius float64, dead bool) {
	p.hp, p.radius, p.dead = hp, radius, dead
}

type Players interface {
	Me() Me
	Add(p Player)
	Remove(id uint64)
	Get(id uint64) (Player, bool)
	GetAll() []Player
}

type players struct {
	me     Me
	byID   map[uint64]Player
	mxByID sync.RWMutex
}

func NewPlayers() Players {
	return &players{
		me:   newEmptyMe(),
		byID: make(map[uint64]Player),
	}
}

func (pp *players) Me() Me { return pp.me }

func (pp *players) Add(p Player) {
	pp.mxByID.Lock()
	defer pp.mxByID.Unlock()
	pp.byID[p.ID()] = p
}

func (pp *players) Remove(id uint64) {
	if pp.me.ID() == id {
		// TODO quit or problems with network
	}
	pp.mxByID.Lock()
	defer pp.mxByID.Unlock()
	delete(pp.byID, id)
}

func (pp *players) Get(id uint64) (Player, bool) {
	if pp.me.ID() == id {
		return pp.me.GetPlayer(), true
	}
	pp.mxByID.RLock()
	defer pp.mxByID.RUnlock()
	p, ok := pp.byID[id]
	return p, ok
}

func (pp *players) GetAll() []Player {
	out := make([]Player, 0, len(pp.byID))
	func() {
		pp.mxByID.RLock()
		defer pp.mxByID.RUnlock()
		for _, p := range pp.byID {
			out = append(out, p)
		}
	}()
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID() < out[j].ID()
	})
	return out
}
