package scene

import (
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/dqso/mincer/client/internal/hud"
	"github.com/dqso/mincer/client/internal/input"
	"github.com/dqso/mincer/client/internal/render"
	"github.com/hajimehoshi/ebiten/v2"
)

type MincerScene struct {
	input  *input.GameInput
	world  entity.World
	render *render.World

	killMessages []hud.KillMessageRender
}

func NewMincerScene(world entity.World) *MincerScene {
	return &MincerScene{
		input: input.NewGameInput(),
		world: world,
	}
}

func (s *MincerScene) Update(state State) error {
	if s.render == nil {
		if s.world.IsLoaded() {
			s.render = render.NewWorld(s.world)
		}
	}
	select {
	case <-state.events.Disconnected():
		state.manager.Go(NewCloseScene())
		return nil
	default:
	}

	s.input.Update()
	if s.input.BeReborn() && s.world.Players().Me().IsDead() {
		state.events.BeReborn()
	}

	state.world.Players().Me().SetDirection(s.input.Direction())
	state.world.Players().Me().SetAttack(s.input.Attack > 0)

	s.killMessages = hud.NewKillMessages(s.world.KillTable().Get())

	s.render.Update()

	return nil
}

func (s *MincerScene) Draw(screen *ebiten.Image) {
	if s.render != nil {
		s.render.Draw(screen)
	}

	hud.DrawStats(screen, s.world.Players().Me())

	hud.DrawFPS(screen)
	hud.DrawPosition(screen, s.world.Players().Me())

	var y float32 = 0.0
	for _, message := range s.killMessages {
		y += float32(message.Draw(screen, float64(y)))
	}
}
