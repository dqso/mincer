package entity

type PlayerStats interface {
	Class() Class
	PhysicalResist() float64
	MagicalResist() float64
	Radius() float32
	Speed() float32
	MaxHP() int32
}

type playerStats struct {
	class          Class
	physicalResist float64
	magicalResist  float64
	radius         float32
	speed          float32
	maxHP          int32
}

func NewPlayerStats(class Class, physicalResist, magicalResist, radius, speed float64, maxHP int32) PlayerStats {
	return &playerStats{
		class:          class,
		physicalResist: physicalResist,
		magicalResist:  magicalResist,
		radius:         float32(radius),
		speed:          float32(speed),
		maxHP:          maxHP,
	}
}

func (s *playerStats) Class() Class {
	return s.class
}

func (s *playerStats) PhysicalResist() float64 {
	return s.physicalResist
}

func (s *playerStats) MagicalResist() float64 {
	return s.magicalResist
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
