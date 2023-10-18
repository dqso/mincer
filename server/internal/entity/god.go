package entity

import "math/rand"

type God interface {
	Int(min, max int) int
	Float(min, max float64) float64
}

type god struct {
	rand *rand.Rand
}

func NewGod(seed int64) God {
	return &god{
		rand: rand.New(rand.NewSource(seed)),
	}
}

func (g *god) Int(min, max int) int {
	return min + rand.Int()%(max-min+1)
}

func (g *god) Float(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

const (
	defaultPlayerHP int64 = 100

	minPlayerRadius float64 = 10.0
	maxPlayerRadius float64 = 50.0

	minPlayerSpeed float64 = 50.0
	maxPlayerSpeed float64 = 100.0

	MaxNorth float64 = -100.0
	MaxEast  float64 = 100.0
	MaxSouth float64 = 100.0
	MaxWest  float64 = -100.0
)
