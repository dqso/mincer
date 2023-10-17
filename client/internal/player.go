package mincer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type Object interface {
	X() float32
	Y() float32
	Size() float32
	Draw(screen *ebiten.Image)
}

type ObjectCircle struct {
	x, y   float32
	radius float32
	color  color.Color
}

func NewObjectCircle(x, y float32, radius float32, color color.Color) *ObjectCircle {
	return &ObjectCircle{
		x:      x,
		y:      y,
		radius: radius,
		color:  color,
	}
}

func (o ObjectCircle) X() float32    { return o.x }
func (o ObjectCircle) Y() float32    { return o.y }
func (o ObjectCircle) Size() float32 { return o.radius * 2 }

func (o ObjectCircle) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, o.x, o.y, o.radius, o.color, true)
}
