package entity

import "math"

type Point struct {
	X float64
	Y float64
}

func (p Point) SubtractPoint(p2 Point) Point {
	return Point{X: p.X - p2.X, Y: p.Y - p2.Y}
}

func (p Point) Distance(p2 Point) float64 {
	s := p2.SubtractPoint(p)
	return math.Sqrt(s.X*s.X + s.Y*s.Y)
}

func (p Point) Middle(p2 Point) Point {
	return Point{
		X: (p.X + p2.X) / 2,
		Y: (p.Y + p2.Y) / 2,
	}
}

type Rect struct {
	LeftUp    Point
	RightDown Point
}
