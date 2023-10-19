package nc_handler

import (
	"context"
	"github.com/dqso/mincer/server/internal/api"
)

func init() { register(api.Code_QUIT, (*Quit)(nil)) }

type Quit struct {
	api.Quit
}

func (r *Quit) Validate() error {
	return nil
}

func (r *Quit) Execute(ctx context.Context, fromClientID uint64, uc usecase) error {
	return uc.Quit(ctx, fromClientID)
}
