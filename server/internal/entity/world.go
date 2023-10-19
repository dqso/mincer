package entity

import "math"

type World interface {
	NewPlayer(id uint64) (Player, error)
	Players() PlayerList
}

type world struct {
	westNorth Point
	eastSouth Point
	players   PlayerList
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

func (w *world) NewPlayer(id uint64) (Player, error) {
	if _, ok := w.players.Get(id); ok {
		return nil, ErrPlayerAlreadyExists
	}
	class := Classes()[w.God().Int(0, len(Classes())-1)]
	p := NewPlayer(id, class)
	radius := p.Radius()
	p.SetPosition(Point{
		X: w.God().Float(radius, w.Width()-radius),
		Y: w.God().Float(radius, w.Height()-radius),
	})
	w.players.Add(p)
	return p, nil
}

func (w *world) Players() PlayerList {
	return w.players
}
