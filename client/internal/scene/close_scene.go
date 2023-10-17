package scene

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type CloseScene struct {
	msg    string
	frames int
}

func NewCloseScene() *CloseScene {
	s := &CloseScene{
		msg: "bye",
	}
	return s
}

func (s *CloseScene) Update(state State) error {
	if s.frames <= 40 {
		s.frames++
		return nil
	}

	return fmt.Errorf("game ended by player")
}

func (s *CloseScene) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Dimgray)
	ebitenutil.DebugPrint(screen, s.msg)
}
