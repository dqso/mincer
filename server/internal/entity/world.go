package entity

import (
	"context"
	"math"
	"slices"
	"sync"
	"time"
)

type World interface {
	Northwest() Point
	Southeast() Point

	NewPlayer(id uint64) (Player, error)
	Respawn(p Player)
	SizeRect() Rect
	Players() PlayerList
	SearchNearby(point Point, callback func(p Player) (stop bool)) Player
	ProjectileList() ProjectileList
	Horn() Horn

	AcquireClass() Class
	AcquireWeapon(class Class) Weapon

	RegisterBot(b Bot, stopFunc func())
	AcquireBotID() uint32
}

type world struct {
	horn Horn

	northwest  Point
	southeast  Point
	playerList PlayerList
	god        God

	botList   BotList
	regions   map[int16][]uint64
	mxRegions sync.RWMutex

	projectileList ProjectileList
}

func NewWorld(ctx context.Context, seed int64, northwest, southeast Point, horn Horn) World {
	w := &world{
		horn:       horn,
		northwest:  northwest,
		southeast:  southeast,
		playerList: NewPlayerList(),
		god:        NewGod(seed),
		botList:    NewBotList(),

		regions: make(map[int16][]uint64),

		projectileList: newProjectileList(),
	}

	go w.supportRegions(ctx)

	return w
}

func (w *world) Northwest() Point {
	return w.northwest
}

func (w *world) Southeast() Point {
	return w.southeast
}

func (w *world) Horn() Horn {
	return w.horn
}

func (w *world) Width() float64 {
	return math.Abs(w.northwest.X - w.southeast.X)
}

func (w *world) Height() float64 {
	return math.Abs(w.northwest.Y - w.southeast.Y)
}

func (w *world) God() God {
	return w.god
}

func (w *world) SizeRect() Rect {
	return Rect{
		LeftUp:    w.northwest,
		RightDown: w.southeast,
	}
}

func (w *world) NewPlayer(id uint64) (Player, error) {
	if _, ok := w.playerList.Get(id); ok {
		return nil, ErrPlayerAlreadyExists
	}
	class := w.AcquireClass()
	weapon := w.AcquireWeapon(class)
	p := NewPlayer(id, class, weapon, w.Horn())
	p.SetPosition(w.acquirePosition(p.Radius()))
	w.playerList.Add(p)
	return p, nil
}

func (w *world) RegisterBot(b Bot, stopFunc func()) {
	w.botList.Add(b, stopFunc)
	b.SetPosition(w.acquirePosition(b.Radius()))
	w.playerList.Add(b)
}

func (w *world) AcquireBotID() uint32 { return w.botList.AcquireID() }

func (w *world) Respawn(p Player) {
	class := w.AcquireClass()
	weapon := w.AcquireWeapon(class)
	p.SetStats(newPlayerStats(class, defaultPlayerRadius, defaultPlayerHP))
	p.SetWeapon(weapon)
	p.SetHP(p.MaxHP())
	p.SetPosition(w.acquirePosition(p.Radius()))
}

func (w *world) AcquireClass() Class {
	list := Classes()
	return list[w.God().Int(0, len(list)-1)]
}

func (w *world) AcquireWeapon(class Class) Weapon {
	list := Weapons(class)
	return list[w.God().Int(0, len(list)-1)]()
}

func (w *world) acquirePosition(radius float64) Point {
	return Point{
		X: w.God().Float(radius, w.Width()-radius),
		Y: w.God().Float(radius, w.Height()-radius),
	}
}

func (w *world) Players() PlayerList {
	return w.playerList
}

const nr int16 = 10 // nr * nr is an amount of regions

func (w *world) supportRegions(ctx context.Context) {
	deltaX := math.Abs(w.southeast.X-w.northwest.X) / float64(nr)
	deltaY := math.Abs(w.southeast.Y-w.northwest.Y) / float64(nr)

	dpx, dpy := 0.0, 0.0
	if w.northwest.X < 0 {
		dpx = math.Abs(w.northwest.X)
	}
	if w.northwest.Y < 0 {
		dpy = math.Abs(w.northwest.Y)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		newRegions := make(map[int16][]uint64)
		for _, p := range w.Players().Slice() {
			point := p.Position()
			regionIdx :=
				int16(math.Floor((point.X+dpx)/deltaX)) +
					int16(math.Floor((point.Y+dpy)/deltaY))*nr
			r, ok := newRegions[regionIdx]
			if !ok {
				newRegions[regionIdx] = []uint64{p.ID()}
			} else {
				newRegions[regionIdx] = append(r, p.ID())
			}
		}

		func() {
			w.mxRegions.Lock()
			defer w.mxRegions.Unlock()
			w.regions = newRegions
		}()

		time.Sleep(time.Millisecond * 500)
	}
}

func (w *world) calculateRegion(p Point) int16 {
	deltaX := math.Abs(w.southeast.X-w.northwest.X) / float64(nr)
	deltaY := math.Abs(w.southeast.Y-w.northwest.Y) / float64(nr)

	dpx, dpy := 0.0, 0.0
	if w.northwest.X < 0 {
		dpx = math.Abs(w.northwest.X)
	}
	if w.northwest.Y < 0 {
		dpy = math.Abs(w.northwest.Y)
	}
	return int16(math.Floor((p.X+dpx)/deltaX)) +
		int16(math.Floor((p.Y+dpy)/deltaY))*nr
}

// bfs
func (w *world) SearchNearby(point Point, cb func(p Player) bool) Player {
	w.mxRegions.RLock()
	defer w.mxRegions.RUnlock()

	region := w.calculateRegion(point)

	seen := make(map[int16]struct{})
	const nr2 = nr * nr

	step := map[int16]struct{}{
		region: {},
	}
	for len(step) > 0 {
		nextStep := make(map[int16]struct{})
		for region := range step {
			if _, ok := seen[region]; ok {
				continue
			}
			for _, id := range func() []uint64 {
				s, ok := w.regions[region]
				if !ok {
					return []uint64{}
				}
				return slices.Clone(s)
			}() {
				p, ok := w.playerList.Get(id)
				if !ok {
					continue
				}
				if cb(p) {
					return p
				}
			}
			seen[region] = struct{}{}
			if region+1 < nr2 { // right
				nextStep[region+1] = struct{}{}
			}
			if region-1 >= 0 { // left
				nextStep[region-1] = struct{}{}
			}
			if region+nr < nr2 { // down
				nextStep[region+nr] = struct{}{}
			}
			if region-nr >= 0 { // up
				nextStep[region-nr] = struct{}{}
			}
		}
		step = nextStep
	}

	return nil
}

func (w *world) ProjectileList() ProjectileList {
	return w.projectileList
}
