package scene

import (
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/dqso/mincer/client/internal/hud"
	"github.com/dqso/mincer/client/internal/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type MincerScene struct {
	input  *input.GameInput
	world  entity.World
	cx, cy float32

	killMessages []hud.KillMessageRender
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

	s.killMessages = hud.NewKillMessages(s.world.KillTable().Get())

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

func (s *MincerScene) drawHUD(screen *ebiten.Image) {
	hud.DrawStats(screen, s.world.Players().Me())

	hud.DrawFPS(screen)

	var y float32 = 0.0
	for _, message := range s.killMessages {
		y += float32(message.Draw(screen, float64(y)))
	}
}
