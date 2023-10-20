package hud

import (
	"fmt"
	"github.com/dqso/mincer/client/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

func DrawFPS(screen *ebiten.Image) {
	str := fmt.Sprintf("%0.2f FPS", ebiten.ActualFPS())
	y := screen.Bounds().Max.Y - fonts.Normal.Metrics().XHeight.Round()
	text.Draw(screen, str, fonts.Normal, 5, y, color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
}
