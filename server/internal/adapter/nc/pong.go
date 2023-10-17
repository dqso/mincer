package nc_adapter

import (
	"context"
	"github.com/dqso/mincer/server/internal/api"
)

func (p *Producer) Pong(ctx context.Context, toClientID uint64, pong string) error {
	bts, err := p.marshalMessage(api.Code_PONG, &api.Pong{
		Pong: pong,
	})
	if err != nil {
		return err
	}
	return p.server.SendPayloadToClient(toClientID, bts)
}
