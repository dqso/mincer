package entity

import "math"

type World interface {
	AddPlayer(id uint64) (Player, error)
	RemovePlayer(id uint64)
	Players() Players
}

type world struct {
	westNorth Point
	eastSouth Point
	players   Players
	god       God
}

func NewWorld(seed int64, westNorth, eastSouth Point) World {
	return &world{
		westNorth: westNorth,
		eastSouth: eastSouth,
		players:   NewPlayers(),
		god:       NewGod(seed),
	}
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

func (w *world) AddPlayer(id uint64) (Player, error) {
	if _, ok := w.players.Get(id); ok {
		return nil, ErrPlayerAlreadyExists
	}
	p := NewPlayer(id)
	radius := p.Radius()
	p.SetPosition(
		w.God().Float(radius, w.Width()-radius),
		w.God().Float(radius, w.Height()-radius),
	)
	w.players.Add(p)
	return p, nil
}

func (w *world) RemovePlayer(id uint64) {
	w.players.Remove(id)
}

func (w *world) Players() Players {
	return w.players
}
