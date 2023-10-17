package scene

import (
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type MincerScene struct {
	world  *entity.World
	cx, cy float32
}

func NewMincerScene(world *entity.World) *MincerScene {
	return &MincerScene{
		world: world,
	}
}

func (s *MincerScene) Update(state State) error {
	s.cx, s.cy = float32(state.world.Me.X), float32(state.world.Me.X)
	return nil
}

func (s *MincerScene) Draw(screen *ebiten.Image) {
	ids := s.world.PlayerIDs()
	for _, id := range ids {
		p, ok := s.world.Player(id)
		if !ok {
			continue
		}
		s.drawPlayer(screen, p)
	}
}

func (s *MincerScene) drawPlayer(screen *ebiten.Image, p *entity.Player) {
	x := float32(p.X) - s.cx + float32(screen.Bounds().Dx())/2
	y := float32(p.Y) - s.cy + float32(screen.Bounds().Dy())/2
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)+1, colornames.White, true)
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius), colornames.Darkred, true)
}
