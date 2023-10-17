package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MincerScene struct{}

func NewMincerScene() *MincerScene {
	return &MincerScene{}
}

func (s *MincerScene) Update(state State) error {
	return nil
}

func (s *MincerScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "mincer")
}
