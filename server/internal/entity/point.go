package entity

type Point struct {
	X float64
	Y float64
}

type Rect struct {
	LeftUp    Point
	RightDown Point
}
