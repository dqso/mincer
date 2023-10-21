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
	"github.com/dqso/mincer/server/internal/usecase/token"
	usecase_world "github.com/dqso/mincer/server/internal/usecase/world"
	"github.com/dqso/mincer/server/pkg/nc"
	"github.com/dqso/mincer/server/pkg/postgres"
	"github.com/dqso/mincer/server/pkg/shutdown"
	"log"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = time.Second * 60

func main() {
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var closer shutdown.Closer
	defer func() {
		log.Println("shutting down server gracefully") // TODO logger
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := closer.Close(ctx); err != nil {
			log.Print(err) // TODO logger
			return
		}
	}()

	config, err := configuration.NewConfig()
	if err != nil {
		log.Print(err) // TODO logger
		return
	}

	pgPool, err := postgres.Connect(ctx, config)
	if err != nil {
		log.Print(err) // TODO logger
		return
	}
	closer.Add(func(ctx context.Context) error {
		pgPool.Close()
		return nil
	})

	ncServer, err := nc.Connect(config)
	if err != nil {
		log.Print(err)
		return
	}
	closer.Add(func(ctx context.Context) error {
		return ncServer.Stop()
	})

	ncProducer := nc_adapter.NewProducer(config, ncServer)
	ncProducerDone := ncProducer.StartLoop(ctx)
	closer.Add(func(ctx context.Context) error {
		select {
		case <-ncProducerDone:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	repoWorld := repository_world.NewRepository(ctx, pgPool)

	usecaseWorld := usecase_world.NewUsecase(ncProducer, repoWorld)

	ncConsumer, err := nc_handler.NewConsumer(ctx, config, ncServer, usecaseWorld)
	if err != nil {
		log.Print(err)
		return
	}
	closer.Add(ncConsumer.Close)
	log.Printf("netcode server started on %s...", config.NCAddress())

	repositoryToken := repository_token.NewRepository(pgPool)

	usecaseToken := usecase_token.NewUsecase(config, repositoryToken)

	handler := rest.NewHandler(usecaseToken)

	httpServer := rest.NewServer(config, handler)
	if err := httpServer.Start(ctx); err != nil {
		log.Print(err) // TODO logger
		return
	}
	closer.Add(httpServer.Close)
	log.Printf("rest server started on %s...", config.RestAddress())

	<-ctx.Done()
}
