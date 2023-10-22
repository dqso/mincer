package hud

import (
	"github.com/dqso/mincer/client/fonts"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"math"
	"strings"
)

type ActionMessageRender struct {
	Alpha uint8
	Words []WordRender
}

type WordRender struct {
	Word   entity.WordRender
	Glyphs []text.Glyph
}

func NewActionMessages(messages []entity.ActionMessage) []ActionMessageRender {
	out := make([]ActionMessageRender, 0, len(messages))

	for _, message := range messages {
		km := ActionMessageRender{
			Alpha: message.Alpha,
		}
		for _, word := range message.Message.Words() {
			km.Words = append(km.Words, WordRender{
				Word:   word,
				Glyphs: text.AppendGlyphs(make([]text.Glyph, 0, len(word.Text())), fonts.Normal, word.Text()),
			})
		}
		out = append(out, km)
	}
	return out
}

func (m *ActionMessageRender) Draw(screen *ebiten.Image, y float64) float64 {
	var str strings.Builder
	for idx, word := range m.Words {
		str.WriteString(word.Word.Text())
		if idx != 0 {
			str.WriteRune(' ')
		}
	}
	bounds := fonts.BoundString(fonts.Normal, str.String())
	x := float64(screen.Bounds().Max.X - bounds.Max.X.Round() - 20)
	y += math.Abs(float64(fonts.Normal.Metrics().Height.Round()))
	glX := 0.0

	op := &ebiten.DrawImageOptions{}
	for _, word := range m.Words {
		c := word.Word.Color()
		c.A = 255 - m.Alpha
		op.ColorScale.Reset()
		op.ColorScale.ScaleWithColor(c)
		for _, gl := range word.Glyphs {
			op.GeoM.Reset()
			op.GeoM.Translate(x, y)
			op.GeoM.Translate(gl.X, gl.Y)
			screen.DrawImage(gl.Image, op)
			glX = gl.X + float64(gl.Image.Bounds().Max.X)
		}
		x += glX + 5
	}

	return float64(fonts.Normal.Metrics().Height.Round())
}
