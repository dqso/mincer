package nc_handler

import (
	"context"
	"fmt"
	"github.com/dqso/mincer/server/internal/api"
)

func init() { register(api.Code_CLIENT_INFO, (*ClientInfo)(nil)) }

type ClientInfo struct {
	api.ClientInfo
}

func (r *ClientInfo) Validate() error {
	if r.Direction < 0 || r.Direction > 360 {
		return fmt.Errorf("direction is not in range [0; 360]")
	}
	if r.DirectionAim < 0 || r.DirectionAim > 360 {
		return fmt.Errorf("directionAim is not in range [0; 360]")
	}
	return nil
}

func (r *ClientInfo) Execute(ctx context.Context, fromClientID uint64, uc usecase) error {
	return uc.ClientInfo(ctx, fromClientID, r.Direction, r.IsMoving, r.Attack, r.DirectionAim)
}
