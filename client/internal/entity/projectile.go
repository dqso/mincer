package entity

import (
	"image/color"
)

type Projectile interface {
	ID() uint64
	Color() color.NRGBA
	Position() Point
	SetPosition(p Point)
	Radius() float64
	Speed() float64
	Direction() float64
}

type projectile struct {
	id        uint64
	color     color.NRGBA
	position  Point
	radius    float64
	speed     float64
	direction float64
}

func NewProjectile(id uint64, color color.NRGBA, position Point, radius, speed, direction float64) Projectile {
	return &projectile{
		id:        id,
		color:     color,
		position:  position,
		radius:    radius,
		speed:     speed,
		direction: direction,
	}
}

func (p *projectile) ID() uint64          { return p.id }
func (p *projectile) Color() color.NRGBA  { return p.color }
func (p *projectile) Position() Point     { return p.position }
func (p *projectile) SetPosition(v Point) { p.position = v }
func (p *projectile) Radius() float64     { return p.radius }
func (p *projectile) Speed() float64      { return p.speed }
func (p *projectile) Direction() float64  { return p.direction }
