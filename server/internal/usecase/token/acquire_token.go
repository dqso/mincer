package usecase_token

import (
	"context"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/pkg/nc"
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

	ip, err := nc.GetLocalIP()
	if err != nil {
		return 0, nil, err
	}

	token := netcode.NewConnectToken()
	err = token.Generate(clientID,
		[]net.UDPAddr{{IP: ip, Port: uc.config.NCPort()}},
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
