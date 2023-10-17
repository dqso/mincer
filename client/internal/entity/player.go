package entity

type Player struct {
	ID     uint64
	X, Y   float64
	HP     int64
	Radius float64
	Dead   bool
}

type Me struct {
	*Player
	Direction float64
	Speed     float64
}
