package entity

type World interface {
	AddPlayer(id uint64) (Player, error)
	Players() Players
}

type world struct {
	width   float64
	height  float64
	players Players
	god     God
}

func NewWorld(seed int64, westNorth, eastSouth Point) World {
	return &world{
		players: NewPlayers(),
		god:     NewGod(seed),
	}
}

func (w *world) Width() float64 {
	return w.width
}

func (w *world) Height() float64 {
	return w.height
}

func (w *world) God() God {
	return w.god
}

func (w *world) AddPlayer(id uint64) (Player, error) {
	if w.players.IsExists(id) {
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

func (w *world) Players() Players {
	return w.players
}
