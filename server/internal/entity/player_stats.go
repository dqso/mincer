package entity

import "sync"

type PlayerStats interface {
	Class() Class
	SetClass(v Class)

	Radius() float64
	SetRadius(v float64)

	Speed() float64
	SetSpeed(v float64)

	MaxHP() int32
	SetMaxHP(v int32)
}

type playerStats struct {
	mx sync.RWMutex

	class  Class
	radius float64
	speed  float64
	maxHP  int32
}

func newPlayerStats(class Class, radius, speed float64, maxHP int32) PlayerStats {
	return &playerStats{
		class:  class,
		radius: radius,
		speed:  speed,
		maxHP:  maxHP,
	}
}

func (s *playerStats) Class() Class {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.class
}

func (s *playerStats) SetClass(v Class) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.class = v
}

func (s *playerStats) Radius() float64 {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.radius
}

func (s *playerStats) SetRadius(v float64) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.radius = v
}

func (s *playerStats) Speed() float64 {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.speed
}

func (s *playerStats) SetSpeed(v float64) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.speed = v
}

func (s *playerStats) MaxHP() int32 {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.maxHP
}

func (s *playerStats) SetMaxHP(v int32) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.maxHP = v
}
