package network

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/dqso/mincer/client/internal/api"
	"github.com/wirepair/netcode"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"time"
)

type NetcodeToken struct {
	ClientID     uint64 `json:"client_id"`
	ConnectToken string `json:"connect_token"`
}

type Manager struct {
	tokenUrl string
	nc       *netcode.Client

	onConnected chan struct{}
}

func NewManager(tokenUrl string) *Manager {
	m := &Manager{
		tokenUrl:    tokenUrl,
		onConnected: make(chan struct{}),
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
		m.nc = netcode.NewClient(token)
		m.nc.SetId(id)
		if err := m.nc.Connect(); err != nil {
			log.Print(err)
			return
		}

		clientTime := float64(0)
		delta := float64(1.0 / 60.0)
		deltaTime := time.Duration(delta * float64(time.Second))
		qwe := false
		for {
			if clientTime > 1 && !qwe {
				qwe = true
				log.Print("ping send")
				bts, err := proto.Marshal(&api.Ping{Ping: "ping"})
				if err != nil {
					log.Print(err)
					continue
				}
				bts, err = proto.Marshal(&api.Message{Code: api.Code_PING, Payload: bts})
				if err != nil {
					log.Print(err)
					continue
				}
				if err := m.nc.SendData(bts); err != nil {
					log.Print(err)
					continue
				}
			}

			m.nc.Update(clientTime)

			data, n := m.nc.RecvData()
			if n > 0 {
				log.Printf("%d: %s", n, string(data))
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

func (m *Manager) getConnectToken(ctx context.Context) (uint64, *netcode.ConnectToken, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.tokenUrl, bytes.NewReader([]byte("{}")))
	if err != nil {
		return 0, nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	var response NetcodeToken
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, nil, err
	}
	log.Printf("%d: %s", response.ClientID, response.ConnectToken)

	tokenBts, err := base64.StdEncoding.DecodeString(response.ConnectToken)
	if err != nil {
		return 0, nil, err
	}
	connToken, err := netcode.ReadConnectToken(tokenBts)
	if err != nil {
		return 0, nil, err
	}
	return response.ClientID, connToken, nil
}
