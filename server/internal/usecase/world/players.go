package usecase_world

type Player struct {
	clientID uint64
	x, y     float64
	size     float64

	direction float64
	speed     float64
}

type Players map[uint64]Player
