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

	actionMessages []hud.ActionMessageRender
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

	s.actionMessages = hud.NewActionMessages(s.world.ActionTable().Get())

	s.render.Update()

	state.world.Players().Me().SetAttack(s.input.Attack > 0, s.input.MouseDirection(s.render.MePosition))

	return nil
}

func (s *MincerScene) Draw(screen *ebiten.Image) {
	if s.render != nil {
		s.render.Draw(screen)
	}

	me := s.world.Players().Me()

	if me.IsLoaded() {
		hud.DrawStats(screen, me)
		hud.DrawClassAndWeapon(screen, me.Class(), me.Weapon())
	}

	hud.DrawFPS(screen)
	if me.IsLoaded() {
		hud.DrawPosition(screen, me)
	}

	var y float32 = 0.0
	for _, message := range s.actionMessages {
		y += float32(message.Draw(screen, float64(y)))
	}
}
