package scene

import (
	"fmt"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/dqso/mincer/client/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image/color"
)

type MincerScene struct {
	input  *input.GameInput
	world  entity.World
	cx, cy float32
}

func NewMincerScene(world entity.World) *MincerScene {
	return &MincerScene{
		input: input.NewGameInput(),
		world: world,
	}
}

func (s *MincerScene) Update(state State) error {
	select {
	case <-state.events.Disconnected():
		state.manager.Go(NewCloseScene())
		return nil
	default:
	}
	s.input.Update()
	state.world.Players().Me().SetDirection(s.input.Direction())
	state.world.Players().Me().SetAttack(s.input.Attack > 0)
	s.cx, s.cy = s.world.Players().Me().Position()
	return nil
}

func (s *MincerScene) Draw(screen *ebiten.Image) {
	for _, player := range s.world.Players().GetAll() {
		s.drawPlayer(screen, player, nil)
	}
	if me := s.world.Players().Me(); me != nil {
		s.drawPlayer(screen, s.world.Players().Me(), entity.ColorBorderMe())
	}
	s.drawHUD(screen)
}

const playerBorderRadius = 1.5

func (s *MincerScene) drawPlayer(screen *ebiten.Image, p entity.Player, border color.Color) {
	if !p.IsLoaded() {
		return
	}
	x, y := p.Position()
	x = x - s.cx + float32(screen.Bounds().Dx())/2
	y = y - s.cy + float32(screen.Bounds().Dy())/2
	radius := p.Radius()
	if border == nil {
		radius += playerBorderRadius
	} else {
		vector.DrawFilledCircle(screen, x, y, radius+playerBorderRadius, border, true)
	}
	bodyColor := p.Color()
	if p.IsDead() {
		bodyColor = entity.ColorDeadPlayer()
	}
	vector.DrawFilledCircle(screen, x, y, radius, bodyColor, true)
}

const hudElementWidth = 100.0
const hudElementHeight = 15.0

func (s *MincerScene) drawHUD(screen *ebiten.Image) {
	me := s.world.Players().Me()
	if !me.IsLoaded() {
		return
	}
	var x, y float32 = 5, 5

	{
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS %0.2f", ebiten.ActualFPS()), int(x), int(y))
		y += hudElementHeight
	}
	{
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("(%0.2f, %0.2f)", s.cx, s.cy), int(x), int(y))
		y += hudElementHeight
	}

	vector.StrokeRect(screen, x, y, hudElementWidth, hudElementHeight, 2, colornames.Red, true)
	width := hudElementWidth * (float32(me.HP()) / 100.0)
	vector.DrawFilledRect(screen, x, y, width, hudElementHeight, colornames.Red, true)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d", me.HP()), int(x+2), int(y))
	y += hudElementHeight + 5

	vector.StrokeRect(screen, x, y, hudElementWidth, hudElementHeight, 2, colornames.Darkcyan, true)
	current := me.CurrentCoolDown()
	width = hudElementWidth * (current / me.MaxCoolDown())
	vector.DrawFilledRect(screen, x, y, width, hudElementHeight, colornames.Darkcyan, true)
	ebitenutil.DebugPrintAt(screen, "cool down", int(x+2), int(y))
}
