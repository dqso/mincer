package game

import (
	"context"
	"fmt"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	ctx            context.Context
	sceneManager   sceneManager
	networkManager networkManager
	world          entity.World
}

func New(ctx context.Context, sceneManager sceneManager, networkManager networkManager, world entity.World) *Game {
	return &Game{
		ctx:            ctx,
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
	select {
	case <-g.ctx.Done():
		return fmt.Errorf("game ended by OS signal")
	default:
	}
	return g.sceneManager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth int, screenHeight int) {
	return outsideWidth, outsideHeight
}
