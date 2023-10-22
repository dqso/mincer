package rest

import (
	"context"
	"github.com/dqso/mincer/server/internal/log"
	"net/http"
	"time"
)

type Server struct {
	logger log.Logger
	config config
	server *http.Server
}

type config interface {
	RestAddress() string
}

func NewServer(logger log.Logger, config config, handler http.Handler) *Server {
	s := &Server{
		logger: logger.With(log.Module("rest_server")),
		config: config,
		server: &http.Server{
			Addr:    config.RestAddress(),
			Handler: handler,
		},
	}
	return s
}

func (s Server) Start(ctx context.Context) error {
	chErr := make(chan error)
	go func() {
		defer close(chErr)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Error("listen and serve error", log.Err(err))
			chErr <- err
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chErr:
		return err
	case <-time.After(time.Millisecond * 100):
		return nil
	}
}

func (s Server) Close(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
