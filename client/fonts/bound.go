package fonts

import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image/color"
	"strings"
)

func BoundString(face font.Face, str string) fixed.Rectangle26_6 {
	str = strings.TrimRight(str, "\n")
	lines := strings.Split(str, "\n")
	if len(lines) == 0 {
		return fixed.Rectangle26_6{}
	}

	minX := fixed.I(0)
	maxX := fixed.I(0)
	for _, line := range lines {
		a := font.MeasureString(face, line)
		if maxX < a {
			maxX = a
		}
	}

	m := face.Metrics()
	minY := -m.Ascent
	maxY := fixed.Int26_6(len(lines)-1)*m.Height + m.Descent
	return fixed.Rectangle26_6{Min: fixed.Point26_6{X: minX, Y: minY}, Max: fixed.Point26_6{X: maxX, Y: maxY}}
}

func ColorNrgba(c color.RGBA) color.NRGBA {
	return color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: c.A,
	}
}
