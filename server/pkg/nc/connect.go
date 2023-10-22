package nc

import (
	"fmt"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/wirepair/netcode"
	"net"
)

type config interface {
	NCPort() int
	NCPrivateKey() []byte
	NCMaxClients() int
}

func Connect(config config) (*netcode.Server, error) {
	ip, err := GetLocalIP()
	if err != nil {
		return nil, err
	}
	udpAddr, err := net.ResolveUDPAddr("", fmt.Sprintf("%s:%d", ip, config.NCPort()))
	if err != nil {
		return nil, err
	}
	s := netcode.NewServer(udpAddr, config.NCPrivateKey(), entity.NCProtocol, config.NCMaxClients())
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s, nil
}
