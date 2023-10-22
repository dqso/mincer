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

	world entity.World

	connected             chan struct{}
	connectingInformation chan string
	disconnected          chan struct{}
	mustDisconnect        chan struct{}
	beReborn              chan struct{}
}

func NewManager(tokenUrl string, world entity.World) *Manager {
	m := &Manager{
		tokenUrl: tokenUrl,
		world:    world,

		connected:             make(chan struct{}),
		connectingInformation: make(chan string, 10),
		disconnected:          make(chan struct{}),
		mustDisconnect:        make(chan struct{}),
		beReborn:              make(chan struct{}),
	}

	m.start()

	return m
}

func (m *Manager) Connected() chan struct{} {
	return m.connected
}

func (m *Manager) ConnectingInformation() chan string {
	return m.connectingInformation
}

func (m *Manager) Disconnected() chan struct{} {
	return m.disconnected
}

func (m *Manager) MustDisconnect() {
	close(m.mustDisconnect)
}

func (m *Manager) BeReborn() {
	m.beReborn <- struct{}{}
}

func (m *Manager) start() {
	//done := make(chan struct{})
	go func() {
		defer close(m.disconnected)
		id, token, err := m.getConnectToken(context.TODO())
		if err != nil {
			log.Print(err)
			m.connectingInformation <- err.Error()
			return
		}
		log.Printf("server addresses: %v", token.ServerAddrs)
		m.world.Players().Me().SetID(id)
		m.nc = netcode.NewClient(token)
		m.nc.SetId(id)
		if err := m.nc.Connect(); err != nil {
			log.Print(err)
			m.connectingInformation <- err.Error()
			return
		}
		close(m.connected)

		clientTime := float64(0)
		delta := float64(1.0 / 60.0)
		deltaTime := time.Duration(delta * float64(time.Second))
		for {
			if m.nc.GetState() <= netcode.StateDisconnected {
				return
			}
			select {
			case <-m.mustDisconnect:
				if err := m.disconnect(); err != nil {
					log.Print(err) // TODO logger
				}
				log.Print("disconnect...")
				return
			case <-m.beReborn:
				if err := m.sendBeReborn(); err != nil {
					log.Print(err) // TODO logger
				}
			default:
			}
			m.nc.Update(clientTime)

			data, _ := m.nc.RecvData()
			if len(data) > 0 {
				m.decodeMessage(data)
			}

			time.Sleep(deltaTime)
			clientTime += deltaTime.Seconds()

			if err := m.repeatingMessageSend(); err != nil {
				log.Print(err) // TODO logger
				//return
			}
		}
	}()
	//return done
}
