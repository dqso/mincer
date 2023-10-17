package usecase_token

import (
	"context"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/wirepair/netcode"
	"net"
)

func (uc Usecase) AcquireToken(ctx context.Context) (uint64, []byte, error) {
	clientID, err := uc.repository.AcquireClientID(ctx)
	if err != nil {
		return 0, nil, err
	}

	userData, err := netcode.RandomBytes(netcode.USER_DATA_BYTES)
	if err != nil {
		return 0, nil, err
	}

	token := netcode.NewConnectToken()
	err = token.Generate(clientID,
		[]net.UDPAddr{{IP: net.ParseIP("192.168.0.17"), Port: 12345}}, // TODO create table, caching...
		netcode.VERSION_INFO,
		entity.NCProtocol,
		10,
		1,
		0,
		userData,
		uc.config.NCPrivateKey(),
	)
	if err != nil {
		return 0, nil, err
	}
	tokenBts, err := token.Write()
	if err != nil {
		return 0, nil, err
	}
	return clientID, tokenBts, nil
}
