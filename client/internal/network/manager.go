package network

import (
	"context"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/wirepair/netcode"
	"log"
	"time"
)

type Manager struct {
	tokenUrl string
	nc       *netcode.Client
	world    entity.World

	onConnected chan struct{}
}

func NewManager(tokenUrl string, world entity.World) *Manager {
	m := &Manager{
		tokenUrl:    tokenUrl,
		onConnected: make(chan struct{}),
		world:       world,
	}

	m.start()

	return m
}

func (m *Manager) OnConnected() chan struct{} {
	return m.onConnected
}

func (m *Manager) start() {
	//done := make(chan struct{})
	go func() {
		//defer close(done)
		id, token, err := m.getConnectToken(context.TODO())
		if err != nil {
			log.Print(err)
			return
		}
		m.world.Players().Me().SetID(id)
		m.nc = netcode.NewClient(token)
		m.nc.SetId(id)
		if err := m.nc.Connect(); err != nil {
			log.Print(err)
			return
		}
		close(m.onConnected)

		clientTime := float64(0)
		delta := float64(1.0 / 60.0)
		deltaTime := time.Duration(delta * float64(time.Second))
		for {
			m.nc.Update(clientTime)

			data, _ := m.nc.RecvData()
			if len(data) > 0 {
				m.decodeMessage(data)
			}

			time.Sleep(deltaTime)
			clientTime += deltaTime.Seconds()

			//case <-ctx.Done():
			//	log.Printf("network manager closed with error: %v", ctx.Err())
			//	return
			//case <-m.stop:
			//	log.Printf("network manager successfully closed")
			//	return
			//case <-time.After(time.Second):
			//	log.Printf("network manager: tick")
			//	if err := m.nc.SendData([]byte("tick")); err != nil {
			//		log.Printf("tick send error: %v", err)
			//	}
			//}
		}
	}()
	//return done
}
