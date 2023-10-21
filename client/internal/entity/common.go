package entity

import "image"

func PointFromImagePoint(p image.Point) Point {
	return Point{
		X: float64(p.X),
		Y: float64(p.Y),
	}
}

type Point struct {
	X, Y float64
}

func (p Point) Add(x, y float64) Point {
	return Point{
		X: p.X + x,
		Y: p.Y + y,
	}
}
