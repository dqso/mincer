package nc_handler

import (
	"context"
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/log"
	"github.com/wirepair/netcode"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"reflect"
	"time"
)

type Consumer struct {
	logger  log.Logger
	config  config
	server  *netcode.Server
	usecase usecase
	closed  chan struct{}
}

type config interface {
	NCRequestPerSecond() int
	NCMaxClients() int
}

type usecase interface {
	AddBot()

	ClientInfo(ctx context.Context, fromUserID uint64, direction float64, isMoving, attack bool, directionAim float64) error
	Quit(ctx context.Context, fromUserID uint64) error
	BeReborn(ctx context.Context, fromUserID uint64) error
	OnPlayerConnect(connect chan uint64, disconnect chan uint64)
	StartLifeCycle(ctx context.Context) chan struct{}
}

func NewConsumer(ctx context.Context, logger log.Logger, config config, server *netcode.Server, usecase usecase) (*Consumer, error) {
	c := &Consumer{
		logger:  logger.With(log.Module("nc_consumer")),
		config:  config,
		server:  server,
		usecase: usecase,
	}
	if err := c.server.Listen(); err != nil {
		return nil, err
	}
	go func() {
		c.closed = make(chan struct{})
		defer close(c.closed)
		c.listen(ctx)
	}()
	return c, nil
}

func (c *Consumer) listen(ctx context.Context) {
	var serverTime float64
	deltaTime := time.Duration(float64(time.Second) / float64(c.config.NCRequestPerSecond()))
	processedPlayersID := make(chan []uint64, 1)

	go c.usecase.OnPlayerConnect(onPlayerConnectDetector(ctx, processedPlayersID))

	stopped := c.usecase.StartLifeCycle(ctx)

	c.usecase.AddBot()
	c.usecase.AddBot()
	c.usecase.AddBot()
	c.usecase.AddBot()
	c.usecase.AddBot()

	for {
		// TODO startTime := time.Now()
		select {
		case <-ctx.Done():
			<-stopped
			c.logger.Debug("netcode consumer was finished")
			return
		default:
		}
		if err := c.server.Update(serverTime); err != nil {
			c.logger.Error("unable to update the netcode server",
				log.Err(err),
			)
			return
		}
		connectedClients := c.server.GetConnectedClientIds()
		processedPlayersID <- connectedClients
		for _, clientID := range connectedClients {
			idx, err := c.server.GetClientIndexByClientId(clientID)
			if err != nil {
				continue
			}
			bts, _ := c.server.RecvPayload(idx)
			if len(bts) == 0 {
				continue
			}
			c.handleMessage(ctx, clientID, bts)
		}
		time.Sleep(deltaTime)             // TODO sleep deltaTime-(time.Now()-startTime)
		serverTime += deltaTime.Seconds() // TODO add deltaTime
	}
}

func (c *Consumer) handleMessage(ctx context.Context, clientID uint64, bts []byte) {
	logger := c.logger.With(slog.Uint64("client_id", clientID))
	var message api.Message
	if err := proto.Unmarshal(bts, &message); err != nil {
		logger.Error("unable to unmarshal the message",
			log.Err(err),
		)
		return
	}
	logger = logger.With(slog.String("code", message.Code.String()))
	rm, ok := registeredMessages[message.Code]
	if !ok {
		logger.Error("this message code is not registered")
		return
	}
	msg := reflect.New(rm.messageType).Interface().(messageInterface)
	if err := proto.Unmarshal(message.Payload, msg); err != nil {
		logger.Error("unable to unmarshal the message",
			log.Err(err),
		)
		return
	}
	if err := msg.Validate(); err != nil {
		logger.Error("incoming message is invalid",
			log.Err(err),
		)
		return
	}
	if err := msg.Execute(ctx, clientID, c.usecase); err != nil {
		logger.Error("execute error",
			log.Err(err),
		)
		return
	}
}

func (c *Consumer) Close(ctx context.Context) error {
	select {
	case <-c.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type messageInterface interface {
	proto.Message
	Validate() error
	Execute(ctx context.Context, fromClientID uint64, uc usecase) error
}

func register(code api.Code, messageNil messageInterface) {
	rm := registeredMessage{
		code:        code,
		messageType: reflect.TypeOf(messageNil).Elem(),
	}
	registeredMessages[code] = rm
}

type registeredMessage struct {
	code        api.Code
	messageType reflect.Type
}

var registeredMessages = make(map[api.Code]registeredMessage)
