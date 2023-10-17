package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type LoadingScene struct {
	connected chan struct{}
}

func NewLoadingScene(connected chan struct{}) *LoadingScene {
	return &LoadingScene{
		connected: connected,
	}
}

func (s *LoadingScene) Update(state State) error {
	select {
	default:
		return nil
	case <-s.connected:
	}
	state.manager.Go(NewMincerScene())
	return nil
}

func (s *LoadingScene) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Darkgreen)
	ebitenutil.DebugPrint(screen, "connecting to server...")
}
