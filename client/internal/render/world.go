package render

import (
	"fmt"
	"github.com/dqso/mincer/client/fonts"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
	"math"
)

type World struct {
	data   entity.World
	width  float64
	height float64
	img    *ebiten.Image

	worldDelta entity.Point
	screenSize entity.Point
	MePosition entity.Point
}

func NewWorld(world entity.World) *World {
	x1, y1, x2, y2 := world.Size()
	w := &World{
		data:   world,
		width:  math.Abs(x2 - x1),
		height: math.Abs(y2 - y1),
	}
	w.img = ebiten.NewImage(
		int(w.width),
		int(w.height),
	)
	w.img.Fill(colornames.GreenA400)
	for i := float64(0); i < w.width; i += 100.0 {
		text.Draw(w.img, fmt.Sprintf("%0.f", i), fonts.Normal, int(i), 20, colornames.Black)
	}
	return w
}

const (
	borderOut = 50.0
)

func (w *World) Update() {
	me := w.data.Players().Me().Position()

	w.worldDelta = entity.Point{
		X: -me.X + w.screenSize.X/2,
		Y: -me.Y + w.screenSize.Y/2,
	}
	// игрок находится рядом с левой или правой границами мира
	if me.X < w.screenSize.X/2-borderOut {
		w.worldDelta.X = borderOut
	} else if w.width-w.screenSize.X/2+borderOut < me.X {
		w.worldDelta.X = -w.width + w.screenSize.X - borderOut
	}
	if w.width < w.screenSize.X-2*borderOut {
		w.worldDelta.X = w.screenSize.X/2 - w.width/2
	}
	// игрок находится рядом с верхней или нижней границами мира
	if me.Y < w.screenSize.Y/2-borderOut {
		w.worldDelta.Y = borderOut
	} else if w.height-w.screenSize.Y/2+borderOut < me.Y {
		w.worldDelta.Y = -w.height + w.screenSize.Y - borderOut
	}
	if w.height < w.screenSize.Y-2*borderOut {
		w.worldDelta.Y = w.screenSize.Y/2 - w.height/2
	}

	w.MePosition = me
	w.MePosition = w.MePosition.Add(w.worldDelta.X, w.worldDelta.Y)
}

func (w *World) Draw(screen *ebiten.Image) {
	w.screenSize = entity.PointFromImagePoint(screen.Bounds().Max)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.worldDelta.X, w.worldDelta.Y)
	screen.DrawImage(w.img, op)

	for _, projectile := range w.data.ProjectileList().GetAll() {
		w.drawProjectile(screen, projectile)
	}

	for _, player := range w.data.Players().GetAll() {
		w.drawPlayer(screen, player, nil)
	}

	w.drawPlayer(screen, w.data.Players().Me(), entity.ColorBorderMe())
}

const playerBorderRadius = 1.5

func (w *World) drawPlayer(screen *ebiten.Image, player entity.Player, border color.Color) {
	if !player.IsLoaded() {
		return
	}
	pos := player.Position()
	pos = pos.Add(w.worldDelta.X, w.worldDelta.Y)
	radius := player.Radius()
	if border == nil {
		radius += playerBorderRadius
	} else {
		vector.DrawFilledCircle(screen, float32(pos.X), float32(pos.Y), radius+playerBorderRadius, border, true)
	}
	bodyColor := player.Color()
	if player.IsDead() {
		bodyColor = entity.ColorDeadPlayer()
	}
	vector.DrawFilledCircle(screen, float32(pos.X), float32(pos.Y), radius, bodyColor, true)
}

func (w *World) drawProjectile(screen *ebiten.Image, projectile entity.Projectile) {
	pos := projectile.Position()
	pos = pos.Add(w.worldDelta.X, w.worldDelta.Y)
	vector.DrawFilledCircle(screen, float32(pos.X), float32(pos.Y), float32(projectile.Radius()), projectile.Color(), true)
}
