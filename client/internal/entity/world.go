package entity

type World interface {
	Players() Players
}

type world struct {
	players Players
}

func NewWorld() World {
	return &world{
		players: NewPlayers(),
	}
}

func (w *world) Players() Players { return w.players }
