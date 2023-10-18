package scene

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type LoadingScene struct {
	log string // TODO use ring
}

func NewLoadingScene() *LoadingScene {
	return &LoadingScene{}
}

func (s *LoadingScene) Update(state State) error {
	select {
	case <-state.events.Connected():
		state.manager.Go(NewMincerScene(state.world))
	case info := <-state.events.ConnectingInformation():
		s.log = info
	default:
	}
	return nil
}

func (s *LoadingScene) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Darkgreen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("connecting to server...\n%s", s.log))
}
