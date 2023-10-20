package entity

import "math/rand"

type God interface {
	Int(min, max int) int
	Float(min, max float64) float64
}

type Rand interface {
	Int() int
	Float64() float64
}

type god struct {
	rand Rand
}

func NewGod(seed int64) God {
	return &god{
		rand: rand.New(rand.NewSource(seed)),
	}
}

func (g *god) Int(min, max int) int {
	return min + g.rand.Int()%(max-min+1)
}

func (g *god) Float(min, max float64) float64 {
	return min + g.rand.Float64()*(max-min)
}

const (
	defaultPlayerHP int64 = 100

	defaultPlayerRadius   float64 = 10.0
	defaultPlayerSpeed    float64 = 100.0 // per second
	defaultPlayerCoolDown float64 = 0.5   // seconds
	defaultPlayerPower    float64 = 100.0
	DefaultAttackRadius   float64 = 40.0

	MaxNorth float64 = 0.0
	MaxWest  float64 = 0.0
	MaxSouth float64 = 300.0
	MaxEast  float64 = 300.0
)
