package nc_handler

import (
	"context"
	"github.com/dqso/mincer/server/internal/api"
)

func init() { register(api.Code_PING, (*PingMessage)(nil)) }

type PingMessage struct {
	api.Ping
}

func (r *PingMessage) Validate() error {
	return nil
}

func (r *PingMessage) Execute(ctx context.Context, fromClientID uint64, uc usecase) error {
	return uc.Ping(ctx, fromClientID, r.Ping.Ping)
}
