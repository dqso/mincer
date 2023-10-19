package entity

import (
	"context"
	"math"
	"slices"
	"sync"
	"time"
)

type World interface {
	NewPlayer(id uint64) (Player, error)
	NewBot() (Bot, error)
	SizeRect() Rect
	Players() PlayerList
	SearchNearby(point Point, cb func(p Player) Player) Player
}

type world struct {
	westNorth  Point
	eastSouth  Point
	playerList PlayerList
	god        God

	botList   BotList
	regions   map[int16][]uint64
	mxRegions sync.RWMutex
}

func NewWorld(seed int64, westNorth, eastSouth Point) World {
	w := &world{
		westNorth:  westNorth,
		eastSouth:  eastSouth,
		playerList: NewPlayerList(),
		god:        NewGod(seed),
		botList:    NewBotList(),

		regions: make(map[int16][]uint64),
	}

	go w.supportRegions(context.TODO()) // tODO

	return w
}

func (w *world) Width() float64 {
	return math.Abs(w.westNorth.X - w.eastSouth.X)
}

func (w *world) Height() float64 {
	return math.Abs(w.westNorth.Y - w.eastSouth.Y)
}

func (w *world) God() God {
	return w.god
}

func (w *world) SizeRect() Rect {
	return Rect{
		LeftUp:    w.westNorth,
		RightDown: w.eastSouth,
	}
}

func (w *world) NewPlayer(id uint64) (Player, error) {
	if _, ok := w.playerList.Get(id); ok {
		return nil, ErrPlayerAlreadyExists
	}
	p := NewPlayer(id, w.acquireClass())
	p.SetPosition(w.acquirePosition(p.Radius()))
	w.playerList.Add(p)
	return p, nil
}

func (w *world) NewBot() (Bot, error) {
	b := w.botList.NewBot(w, w.acquireClass())
	b.SetPosition(w.acquirePosition(b.Radius()))
	w.playerList.Add(b)
	return b, nil
}

func (w *world) acquireClass() Class {
	return Classes()[w.God().Int(0, len(Classes())-1)]
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

const nr int16 = 10 // amount of regions = nr * nr

func (w *world) supportRegions(ctx context.Context) {
	deltaX := math.Abs(w.eastSouth.X-w.westNorth.X) / float64(nr)
	deltaY := math.Abs(w.eastSouth.Y-w.westNorth.Y) / float64(nr)

	dpx, dpy := 0.0, 0.0
	if w.westNorth.X < 0 {
		dpx = math.Abs(w.westNorth.X)
	}
	if w.westNorth.Y < 0 {
		dpy = math.Abs(w.westNorth.Y)
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

		//log.Print(newRegions)

		time.Sleep(time.Second)
	}
}

func (w *world) calculateRegion(p Point) int16 {
	deltaX := math.Abs(w.eastSouth.X-w.westNorth.X) / float64(nr)
	deltaY := math.Abs(w.eastSouth.Y-w.westNorth.Y) / float64(nr)

	dpx, dpy := 0.0, 0.0
	if w.westNorth.X < 0 {
		dpx = math.Abs(w.westNorth.X)
	}
	if w.westNorth.Y < 0 {
		dpy = math.Abs(w.westNorth.Y)
	}
	return int16(math.Floor((p.X+dpx)/deltaX)) +
		int16(math.Floor((p.Y+dpy)/deltaY))*nr
}

func (w *world) SearchNearby(point Point, cb func(p Player) Player) Player {
	region := w.calculateRegion(point)
	return w.searchNearby(region, cb, make(map[int16]struct{}))
}

func (w *world) searchNearby(region int16, cb func(p Player) Player, seen map[int16]struct{}) Player {
	if _, ok := seen[region]; ok {
		return nil
	}
	for _, id := range func() []uint64 {
		w.mxRegions.RLock()
		defer w.mxRegions.RUnlock()
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
		if p = cb(p); p != nil {
			return p
		}
	}
	seen[region] = struct{}{}
	nr2 := nr * nr
	if region+1 < nr2 { // right
		if p := w.searchNearby(region+1, cb, seen); p != nil {
			return p
		}
	}
	if region-1 >= 0 { // left
		if p := w.searchNearby(region-1, cb, seen); p != nil {
			return p
		}
	}
	if region-nr >= 0 { // up
		if p := w.searchNearby(region-nr, cb, seen); p != nil {
			return p
		}
	}
	if region+nr < nr2 { // down
		if p := w.searchNearby(region+nr, cb, seen); p != nil {
			return p
		}
	}
	return nil
}
