package entity

import "sync"

type PlayerStats interface {
	Class() Class
	SetClass(v Class)

	Radius() float64
	SetRadius(v float64)

	Speed() float64
	SetSpeed(v float64)

	MaxHP() int64
	SetMaxHP(v int64)

	MaxCoolDown() float64
	SetMaxCoolDown(v float64)

	Power() float64
	SetPower(v float64)
}

type playerStats struct {
	mx sync.RWMutex

	class       Class
	radius      float64
	speed       float64
	maxHP       int64
	maxCoolDown float64 // per sec
	power       float64
}

func newPlayerStats(class Class, radius, speed float64, maxHP int64, maxCoolDown, power float64) PlayerStats {
	return &playerStats{
		class:       class,
		radius:      radius,
		speed:       speed,
		maxHP:       maxHP,
		maxCoolDown: maxCoolDown,
		power:       power,
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

func (s *playerStats) MaxHP() int64 {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.maxHP
}

func (s *playerStats) SetMaxHP(v int64) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.maxHP = v
}

func (s *playerStats) MaxCoolDown() float64 {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.maxCoolDown
}

func (s *playerStats) SetMaxCoolDown(v float64) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.maxCoolDown = v
}

func (s *playerStats) Power() float64 {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.power
}

func (s *playerStats) SetPower(v float64) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.power = v
}