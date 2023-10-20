package hud

import (
	"fmt"
	"github.com/dqso/mincer/client/fonts"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"math"
)

type KillMessageRender struct {
	Alpha uint8

	PlayerColor  color.NRGBA
	PlayerText   string
	PlayerGlyphs []text.Glyph

	LabelText   string
	LabelGlyphs []text.Glyph

	KillerColor  color.NRGBA
	KillerText   string
	KillerGlyphs []text.Glyph
}

func NewKillMessages(messages []entity.KillMessage) []KillMessageRender {
	out := make([]KillMessageRender, 0, len(messages))

	for _, message := range messages {
		km := KillMessageRender{
			Alpha:       message.Alpha,
			PlayerColor: message.PlayerClass.Color(),
			PlayerText:  fmt.Sprintf("%d", message.PlayerID),
			LabelText:   "kills",
			KillerColor: message.KillerClass.Color(),
			KillerText:  fmt.Sprintf("%d", message.KillerID),
		}

		if message.IsPlayerBot {
			km.PlayerText = "bot " + km.PlayerText
		}
		km.PlayerGlyphs = text.AppendGlyphs(km.PlayerGlyphs, fonts.Normal, km.PlayerText)

		if message.IsKillerBot {
			km.KillerText = " bot " + km.KillerText
		}
		km.KillerGlyphs = text.AppendGlyphs(km.KillerGlyphs, fonts.Normal, km.KillerText)

		km.LabelGlyphs = text.AppendGlyphs(km.LabelGlyphs, fonts.Normal, km.LabelText)
		out = append(out, km)
	}
	return out
}

func (m *KillMessageRender) Draw(screen *ebiten.Image, y float64) float64 {
	str := fmt.Sprintf("%s%s%s", m.KillerText, m.LabelText, m.PlayerText)
	bounds := fonts.BoundString(fonts.Normal, str)
	x := float64(screen.Bounds().Max.X - bounds.Max.X.Round() - 20)
	y += math.Abs(float64(fonts.Normal.Metrics().Height.Round()))
	glX := 0.0

	op := &ebiten.DrawImageOptions{}
	drawColorfulLine := func(c color.NRGBA, glyphs []text.Glyph) {
		c.A = 255 - m.Alpha
		op.ColorScale.Reset()
		op.ColorScale.ScaleWithColor(c)
		for _, gl := range glyphs {
			op.GeoM.Reset()
			op.GeoM.Translate(x, y)
			op.GeoM.Translate(gl.X, gl.Y)
			screen.DrawImage(gl.Image, op)
			glX = gl.X + float64(gl.Image.Bounds().Max.X)
		}
		x += glX + 5
	}

	// Draw the killer
	drawColorfulLine(m.KillerColor, m.KillerGlyphs)

	// Draw the " kills "
	drawColorfulLine(color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF - m.Alpha}, m.LabelGlyphs)

	// Draw the player
	drawColorfulLine(m.PlayerColor, m.PlayerGlyphs)

	return float64(fonts.Normal.Metrics().Height.Round())
}
