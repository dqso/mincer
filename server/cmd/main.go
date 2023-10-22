package main

//go:generate protoc --proto_path=../../proto --go_out=. --go_opt=Mserver_client.proto=./../internal/api server_client.proto

import (
	"context"
	"github.com/dqso/mincer/server/internal/adapter/nc"
	"github.com/dqso/mincer/server/internal/adapter/repository_token"
	"github.com/dqso/mincer/server/internal/adapter/repository_world"
	"github.com/dqso/mincer/server/internal/configuration"
	"github.com/dqso/mincer/server/internal/handler/nc"
	"github.com/dqso/mincer/server/internal/handler/rest"
	"github.com/dqso/mincer/server/internal/log"
	"github.com/dqso/mincer/server/internal/usecase/token"
	usecase_world "github.com/dqso/mincer/server/internal/usecase/world"
	"github.com/dqso/mincer/server/pkg/nc"
	"github.com/dqso/mincer/server/pkg/postgres"
	"github.com/dqso/mincer/server/pkg/shutdown"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = time.Second * 60

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger := log.New()

	var closer shutdown.Closer
	defer func() {
		logger.Info("shutting down gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := closer.Close(ctx); err != nil {
			logger.Error("failed to shutdown", log.Err(err))
			return
		}
	}()

	config, err := configuration.NewConfig()
	if err != nil {
		logger.Error("unable to initialize a configuration", log.Err(err))
		return
	}
	logger = log.NewWithConfig(config)
	slog.SetDefault(logger)

	pgPool, err := postgres.Connect(ctx, config)
	if err != nil {
		logger.Error("unable to connect to postgres", log.Err(err))
		return
	}
	closer.Add(func(ctx context.Context) error {
		pgPool.Close()
		return nil
	})

	ncServer, err := nc.Connect(config)
	if err != nil {
		logger.Error("unable to start a netcode server", log.Err(err))
		return
	}
	closer.Add(func(ctx context.Context) error {
		return ncServer.Stop()
	})

	ncProducer := nc_adapter.NewProducer(config, logger, ncServer)
	ncProducerDone := ncProducer.StartLoop(ctx)
	closer.Add(func(ctx context.Context) error {
		select {
		case <-ncProducerDone:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	repoWorld := repository_world.NewRepository(ctx, logger, pgPool)

	usecaseWorld := usecase_world.NewUsecase(ctx, logger, ncProducer, repoWorld)

	ncConsumer, err := nc_handler.NewConsumer(ctx, logger, config, ncServer, usecaseWorld)
	if err != nil {
		logger.Error("unable to initialize a netcode consumer", log.Err(err))
		return
	}
	closer.Add(ncConsumer.Close)
	logger.Info("netcode server started...", slog.Any("address", config.NCAddress()))

	repositoryToken := repository_token.NewRepository(logger, pgPool)

	usecaseToken := usecase_token.NewUsecase(config, repositoryToken)

	handler := rest.NewHandler(logger, usecaseToken)

	httpServer := rest.NewServer(logger, config, handler)
	if err := httpServer.Start(ctx); err != nil {
		logger.Info("unable to start a rest server", log.Err(err))
		return
	}
	closer.Add(httpServer.Close)
	logger.Info("rest server started...", slog.Any("address", config.RestAddress()))

	<-ctx.Done()
}
