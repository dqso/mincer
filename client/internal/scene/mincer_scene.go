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
	s.cx, s.cy = s.world.Players().Me().PositionFloat32()
	return nil
}

func (s *MincerScene) Draw(screen *ebiten.Image) {
	for _, player := range s.world.Players().GetAll() {
		s.drawPlayer(screen, player, colornames.Darkblue)
	}
	s.drawPlayer(screen, s.world.Players().Me(), colornames.Darkgreen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS %0.2f", ebiten.ActualFPS()))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("(%0.2f, %0.2f)", s.cx, s.cy), 0, 15)
}

func (s *MincerScene) drawPlayer(screen *ebiten.Image, p entity.Player, color color.Color) {
	x, y := p.PositionFloat32()
	x = x - s.cx + float32(screen.Bounds().Dx())/2
	y = y - s.cy + float32(screen.Bounds().Dy())/2
	vector.DrawFilledCircle(screen, x, y, p.RadiusFloat32()+1, colornames.White, true)
	vector.DrawFilledCircle(screen, x, y, p.RadiusFloat32(), color, true)
}
