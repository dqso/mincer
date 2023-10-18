package main

//go:generate protoc --proto_path=../../proto --go_out=. --go_opt=Mserver_client.proto=./../internal/api server_client.proto

import (
	"context"
	"flag"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/dqso/mincer/client/internal/game"
	"github.com/dqso/mincer/client/internal/network"
	"github.com/dqso/mincer/client/internal/scene"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"os/signal"
	"syscall"
)

var tokenUrl string

func init() {
	flag.StringVar(&tokenUrl, "a", "http://localhost:8080/token", "")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	world := entity.NewWorld()

	networkManager := network.NewManager(tokenUrl, world)

	sceneManager := scene.NewManager(scene.NewLoadingScene(), networkManager, world)

	//mincer.NewObjectCircle(100, 100, 13, colornames.Red600)

	g := game.New(ctx, sceneManager, networkManager, world)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("mincer")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
