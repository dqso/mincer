package nc_handler

import (
	"context"
	"github.com/dqso/mincer/server/internal/api"
)

func init() { register(api.Code_BE_REBORN, (*BeReborn)(nil)) }

type BeReborn struct {
	api.BeReborn
}

func (r *BeReborn) Validate() error {
	return nil
}

func (r *BeReborn) Execute(ctx context.Context, fromClientID uint64, uc usecase) error {
	return uc.BeReborn(ctx, fromClientID)
}
