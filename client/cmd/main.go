package main

//go:generate protoc --proto_path=../../proto --go_out=. --go_opt=Mserver_client.proto=./../internal/api server_client.proto

import (
	"flag"
	"github.com/dqso/mincer/client/internal/game"
	"github.com/dqso/mincer/client/internal/network"
	"github.com/dqso/mincer/client/internal/scene"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var tokenUrl string

func init() {
	flag.StringVar(&tokenUrl, "a", "http://localhost:8080/token", "")
}

func main() {
	//ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	//defer cancel()

	networkManager := network.NewManager(tokenUrl)
	sceneManager := scene.NewManager(scene.NewLoadingScene(networkManager.OnConnected()))

	//mincer.NewObjectCircle(100, 100, 13, colornames.Red600)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("mincer")
	if err := ebiten.RunGame(game.New(sceneManager, networkManager)); err != nil {
		log.Fatal(err)
	}
}
