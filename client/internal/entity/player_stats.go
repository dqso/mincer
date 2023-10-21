package entity

type PlayerStats interface {
	Class() Class
	Radius() float32
	Speed() float32
	MaxHP() int32
}

type playerStats struct {
	class  Class
	radius float32
	speed  float32
	maxHP  int32
}

func NewPlayerStats(class Class, radius, speed float64, maxHP int32) PlayerStats {
	return &playerStats{
		class:  class,
		radius: float32(radius),
		speed:  float32(speed),
		maxHP:  maxHP,
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

func (s *playerStats) MaxHP() int32 {
	return s.maxHP
}
