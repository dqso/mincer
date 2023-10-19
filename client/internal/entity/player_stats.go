package entity

type PlayerStats interface {
	Class() Class
	Radius() float32
	Speed() float32
	MaxHP() int64
	MaxCoolDown() float32
	Power() float32
}

type playerStats struct {
	class       Class
	radius      float32
	speed       float32
	maxHP       int64
	maxCoolDown float32 // per sec
	power       float32
}

func NewPlayerStats(class Class, radius, speed float64, maxHP int64, maxCoolDown, power float64) PlayerStats {
	return &playerStats{
		class:       class,
		radius:      float32(radius),
		speed:       float32(speed),
		maxHP:       maxHP,
		maxCoolDown: float32(maxCoolDown),
		power:       float32(power),
	}
}

func (s *playerStats) Class() Class {
	return s.class
}

func (s *playerStats) Radius() float32 {
	return s.radius
}

func (s *playerStats) Speed() float32 {
	return s.speed
}

func (s *playerStats) MaxHP() int64 {
	return s.maxHP
}

func (s *playerStats) MaxCoolDown() float32 {
	return s.maxCoolDown
}

func (s *playerStats) Power() float32 {
	return s.power
}
