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

const (
	hudSeparatorSpace = 5.0
	hudElementWidth   = 130.0
)

func DrawStats(screen *ebiten.Image, me entity.Me) {
	if !me.IsLoaded() {
		return
	}
	var x, y float32 = hudSeparatorSpace, hudSeparatorSpace
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
		y += hudElementHeight + hudSeparatorSpace
	}

	hp := me.HP()
	draw(colornames.Indianred, fonts.Bold, fmt.Sprintf("%d HP", hp), float32(hp), float32(me.MaxHP()))

	draw(colornames.Darkcyan, fonts.Normal, "cool down", float32(me.CurrentCoolDown()), float32(me.Weapon().CoolDown()))
}

func DrawClassAndWeapon(screen *ebiten.Image, class entity.Class, weapon entity.Weapon) {
	widthClass := float64(fonts.BoundString(fonts.Bold, class.Name()).Max.X.Round())
	widthName := float64(fonts.BoundString(fonts.Normal, weapon.Name()).Max.X.Round())
	var physicalLabel, magicalLabel string
	numSeparators := float64(2)
	if weapon.PhysicalDamage() > 0 {
		physicalLabel = fmt.Sprintf("%+d", weapon.PhysicalDamage())
		numSeparators++
	}
	widthPhysical := float64(fonts.BoundString(fonts.Bold, physicalLabel).Max.X.Round())
	if weapon.MagicalDamage() > 0 {
		magicalLabel = fmt.Sprintf("%+d", weapon.MagicalDamage())
		numSeparators++
	}
	widthMagical := float64(fonts.BoundString(fonts.Bold, magicalLabel).Max.X.Round())

	const startX = 3*hudSeparatorSpace + hudElementWidth

	x, y := startX, float32(0.0)
	width := max(
		float32(widthName+widthPhysical+widthMagical+numSeparators*hudSeparatorSpace),
		float32(widthClass+2*numSeparators),
	)
	vector.DrawFilledRect(screen, float32(x)-3, 0,
		width+6,
		float32(fonts.Normal.Metrics().Height.Round()*2+2*hudSeparatorSpace)+3,
		color.NRGBA{R: 0x75, G: 0x97, B: 0xDE, A: 0xFF}, true)
	vector.DrawFilledRect(screen, float32(x), 0,
		width,
		float32(fonts.Normal.Metrics().Height.Round()*2+2*hudSeparatorSpace),
		color.NRGBA{R: 0x19, G: 0x38, B: 0x5C, A: 0xFF}, true)

	x, y = startX+hudSeparatorSpace, float32(fonts.Normal.Metrics().Height.Round())
	text.Draw(screen, class.Name(), fonts.Bold, int(x), int(y), class.Color())

	x, y = startX+hudSeparatorSpace, float32(fonts.Normal.Metrics().Height.Round())*2
	text.Draw(screen, weapon.Name(), fonts.Normal, int(x), int(y), color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
	x += widthName + hudSeparatorSpace
	if weapon.PhysicalDamage() > 0 {
		text.Draw(screen, physicalLabel, fonts.Bold, int(x), int(y), entity.Warrior.Color())
		x += widthPhysical + hudSeparatorSpace
	}
	if weapon.MagicalDamage() > 0 {
		text.Draw(screen, magicalLabel, fonts.Bold, int(x), int(y), entity.Mage.Color())
		x += widthMagical + hudSeparatorSpace
	}
}
