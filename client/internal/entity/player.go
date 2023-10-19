package entity

import (
	"github.com/dqso/mincer/client/internal/api"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
	"sort"
	"sync"
)

type Me interface {
	Player
	GetPlayer() Player
	SetID(id uint64)

	Direction() (float64, bool)
	SetDirection(d float64, isMoving bool)
	Attack() bool
	SetAttack(v bool)
}

type me struct {
	Player

	direction float64
	isMoving  bool
	attack    bool
}

func newEmptyMe() Me {
	return &me{
		Player:    newEmptyPlayer(),
		direction: 0,
	}
}

func (m *me) GetPlayer() Player                     { return m.Player }
func (m *me) SetID(id uint64)                       { m.Player.setID(id) }
func (m *me) Direction() (float64, bool)            { return m.direction, m.isMoving }
func (m *me) SetDirection(d float64, isMoving bool) { m.direction, m.isMoving = d, isMoving }
func (m *me) Attack() bool                          { return m.attack }
func (m *me) SetAttack(v bool)                      { m.attack = v }

type Player interface {
	ID() uint64
	setID(id uint64)
	IsLoaded() bool

	PlayerStats
	SetStats(stats PlayerStats)

	Color() color.Color

	HP() int64
	SetHP(hp int64)
	IsDead() bool

	Position() (float32, float32)
	SetPosition(x, y float64)
}

type player struct {
	id uint64

	PlayerStats

	color color.Color

	hp   int64
	x, y float32
}

func newEmptyPlayer() Player {
	return &player{
		color: colornames.White,
	}
}

func NewPlayer(id uint64, hp int64, playerStats PlayerStats, x, y float64) Player {
	return &player{
		id:          id,
		PlayerStats: playerStats,
		color:       playerStats.Class().color(),
		hp:          hp,
		x:           float32(x),
		y:           float32(y),
	}
}

func (p *player) ID() uint64      { return p.id }
func (p *player) IsLoaded() bool  { return p.PlayerStats != nil }
func (p *player) setID(id uint64) { p.id = id }

func (p *player) SetStats(stats PlayerStats) {
	p.PlayerStats = stats
	p.color = stats.Class().color()
}

func (p *player) Color() color.Color { return p.color }

func (p *player) HP() int64      { return p.hp }
func (p *player) SetHP(hp int64) { p.hp = hp }
func (p *player) IsDead() bool   { return p.hp <= 0 }

func (p *player) Position() (float32, float32) { return p.x, p.y }
func (p *player) SetPosition(x, y float64)     { p.x, p.y = float32(x), float32(y) }

type Class api.Class

const (
	Warrior = Class(api.Class_WARRIOR)
	Mage    = Class(api.Class_MAGE)
	Ranger  = Class(api.Class_RANGER)
)

func (c Class) color() color.Color {
	switch c {
	case Warrior:
		return colornames.Red200
	case Mage:
		return colornames.Blue200
	case Ranger:
		return colornames.Green200
	default:
		return colornames.White
	}
}

func ColorBorderMe() color.Color {
	return colornames.White
}

func ColorDeadPlayer() color.Color {
	return color.RGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0}
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
