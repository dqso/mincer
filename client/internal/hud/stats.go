package hud

import (
	"fmt"
	"github.com/dqso/mincer/client/fonts"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image/color"
)

const hudElementWidth = 130.0

func DrawStats(screen *ebiten.Image, me entity.Me) {
	if !me.IsLoaded() {
		return
	}
	var x, y float32 = 5, 5
	var hudElementHeight = float32(fonts.Normal.Metrics().Height.Round())

	draw := func(clr color.RGBA, face font.Face, str string, current, max float32) {
		vector.StrokeRect(screen, x, y, hudElementWidth, hudElementHeight, 2, clr, true)
		width := hudElementWidth * current / max
		vector.DrawFilledRect(screen, x, y, width, hudElementHeight, clr, true)
		bounds := fonts.BoundString(face, str)
		text.Draw(screen, str, face,
			int(x)+hudElementWidth/2-bounds.Max.X.Round()/2,
			int(y)+2+face.Metrics().CapHeight.Round(),
			color.NRGBA{0xff, 0xff, 0xff, 0xff},
		)
		y += hudElementHeight + 5
	}

	hp := me.HP()
	draw(colornames.Indianred, fonts.Bold, fmt.Sprintf("%d HP", hp), float32(hp), float32(me.MaxHP()))

	draw(colornames.Darkcyan, fonts.Normal, "cool down", me.CurrentCoolDown(), 1000 /*TODO from weapon*/)
}
