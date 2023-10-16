package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

type Game struct{}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("game ended by player")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("mincer")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
