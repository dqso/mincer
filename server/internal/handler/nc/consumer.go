package nc_handler

import (
	"context"
	"github.com/dqso/mincer/server/internal/api"
	"github.com/wirepair/netcode"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
	"time"
)

type Consumer struct {
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
	Ping(ctx context.Context, fromUserID uint64, ping string) error
	OnPlayerConnect(chan uint64)
}

func NewConsumer(ctx context.Context, config config, server *netcode.Server, usecase usecase) (*Consumer, error) {
	c := &Consumer{
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
	processedPlayersID := make(chan uint64, c.config.NCMaxClients())

	go c.usecase.OnPlayerConnect(onPlayerConnectDetector(ctx, processedPlayersID))

	for {
		// TODO startTime := time.Now()
		select {
		case <-ctx.Done():
			log.Print("nc consumer", ctx.Err()) // TODO logger
			return
		default:
		}
		if err := c.server.Update(0.0); err != nil {
			log.Print(err) // TODO logger
			return
		}
		for _, clientID := range c.server.GetConnectedClientIds() {
			idx, err := c.server.GetClientIndexByClientId(clientID)
			if err != nil {
				continue
			}
			processedPlayersID <- clientID
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
	var message api.Message
	if err := proto.Unmarshal(bts, &message); err != nil {
		log.Print(err) // TODO logger
		return
	}
	rm, ok := registeredMessages[message.Code]
	if !ok {
		log.Printf("message %s not registered", message.Code.String()) // TODO logger
		return
	}
	msg := reflect.New(rm.messageType).Interface().(messageInterface)
	if err := proto.Unmarshal(message.Payload, msg); err != nil {
		log.Print(err) // TODO logger
		return
	}
	if err := msg.Validate(); err != nil {
		log.Print(err) // TODO logger
		return
	}
	if err := msg.Execute(ctx, clientID, c.usecase); err != nil {
		log.Print(err) // TODO logger
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