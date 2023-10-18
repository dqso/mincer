package game

import (
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	sceneManager   sceneManager
	networkManager networkManager
	world          entity.World
}

func New(sceneManager sceneManager, networkManager networkManager, world entity.World) *Game {
	return &Game{
		sceneManager:   sceneManager,
		networkManager: networkManager,
		world:          world,
	}
}

type sceneManager interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type networkManager interface {
}

func (g *Game) Update() error {
	return g.sceneManager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}
