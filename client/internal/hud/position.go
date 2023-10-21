package hud

import (
	"fmt"
	"github.com/dqso/mincer/client/fonts"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

func DrawPosition(screen *ebiten.Image, me entity.Me) {
	pos := me.Position()
	y := screen.Bounds().Max.Y - fonts.Normal.Metrics().XHeight.Round() - fonts.Normal.Metrics().Height.Round()
	text.Draw(screen, fmt.Sprintf("(%0.2f,%0.2f)", pos.X, pos.Y), fonts.Normal, 5, y, colornames.White)
}
