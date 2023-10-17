package nc

import (
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/wirepair/netcode"
	"net"
)

type config interface {
	NCAddress() string
	NCPrivateKey() []byte
	NCMaxClients() int
}

func Connect(config config) (*netcode.Server, error) {
	udpAddr, err := net.ResolveUDPAddr("", config.NCAddress())
	if err != nil {
		return nil, err
	}
	s := netcode.NewServer(udpAddr, config.NCPrivateKey(), entity.NCProtocol, config.NCMaxClients())
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s, nil
}
