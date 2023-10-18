package scene

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type CloseScene struct {
}

func NewCloseScene() *CloseScene {
	s := &CloseScene{}
	return s
}

func (s *CloseScene) Update(state State) error {
	select {
	case <-state.events.Disconnected():
		return fmt.Errorf("game ended by player")
	default:
	}
	return nil
}

func (s *CloseScene) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Dimgray)
	ebitenutil.DebugPrint(screen, "bye")
}
