package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type LoadingScene struct {
	onConnected func() chan struct{}
}

func NewLoadingScene(onConnected func() chan struct{}) *LoadingScene {
	return &LoadingScene{
		onConnected: onConnected,
	}
}

func (s *LoadingScene) Update(state State) error {
	select {
	default:
		return nil
	case <-s.onConnected():
	}
	state.manager.Go(NewMincerScene(state.world))
	return nil
}

func (s *LoadingScene) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Darkgreen)
	ebitenutil.DebugPrint(screen, "connecting to server...")
}
