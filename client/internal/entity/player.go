package entity

import (
	"github.com/dqso/mincer/client/internal/api"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
	"sort"
	"sync"
)

const (
	botPrefixID = uint32(0xFFFFFFFF)
)

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

	Position() Point
	SetPosition(Point)
}

type player struct {
	id uint64

	PlayerStats

	color color.Color

	hp       int64
	position Point
}

func newEmptyPlayer() Player {
	return &player{
		color: colornames.White,
	}
}

func NewPlayer(id uint64, hp int64, playerStats PlayerStats, position Point) Player {
	return &player{
		id:          id,
		PlayerStats: playerStats,
		color:       playerStats.Class().Color(),
		hp:          hp,
		position:    position,
	}
}

func (p *player) ID() uint64      { return p.id }
func (p *player) IsLoaded() bool  { return p.PlayerStats != nil }
func (p *player) setID(id uint64) { p.id = id }

func (p *player) SetStats(stats PlayerStats) {
	p.PlayerStats = stats
	p.color = stats.Class().Color()
}

func (p *player) Color() color.Color { return p.color }

func (p *player) HP() int64      { return p.hp }
func (p *player) SetHP(hp int64) { p.hp = hp }
func (p *player) IsDead() bool   { return p.hp <= 0 }

func (p *player) Position() Point         { return p.position }
func (p *player) SetPosition(point Point) { p.position = point }

type Class api.Class

const (
	Warrior = Class(api.Class_WARRIOR)
	Mage    = Class(api.Class_MAGE)
	Ranger  = Class(api.Class_RANGER)
)

func (c Class) Color() color.NRGBA {
	switch c {
	case Warrior:
		return color.NRGBA{0xEF, 0x9A, 0x9A, 0xFF}
	case Mage:
		return color.NRGBA{0x90, 0xCA, 0xF9, 0xFF}
	case Ranger:
		return color.NRGBA{0xA5, 0xD6, 0xA7, 0xFF}
	default:
		return color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
	}
}

func ColorDeadPlayer() color.Color {
	return color.NRGBA{R: 0xAA, G: 0xAA, B: 0xAA, A: 0xFF}
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
