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
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %0.2f", ebiten.ActualFPS()))
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("(%0.2f, %0.2f)", s.cx, s.cy), 0, 15)
}

const playerBorderRadius = 1.5

func (s *MincerScene) drawPlayer(screen *ebiten.Image, p entity.Player, border color.Color) {
	x, y := p.Position()
	x = x - s.cx + float32(screen.Bounds().Dx())/2
	y = y - s.cy + float32(screen.Bounds().Dy())/2
	radius := p.Radius()
	if border == nil {
		radius += playerBorderRadius
	} else {
		vector.DrawFilledCircle(screen, x, y, radius+playerBorderRadius, border, true)
	}
	vector.DrawFilledCircle(screen, x, y, radius, p.Color(), true)
}

const hudElementWidth = 100.0
const hudElementHeight = 15.0

func (s *MincerScene) drawHUD(screen *ebiten.Image) {
	me := s.world.Players().Me()
	vector.StrokeRect(screen, 5, 5, hudElementWidth, hudElementHeight, 2, colornames.Red, true)
	hpWidth := hudElementWidth * (1 - 1/float32(me.HP()))
	vector.DrawFilledRect(screen, 5, 5, hpWidth, hudElementHeight, colornames.Red, true)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d", me.HP()), 7, 6)
}
