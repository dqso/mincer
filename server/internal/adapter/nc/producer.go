package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/wirepair/netcode"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Producer struct {
	config config
	server *netcode.Server

	onPlayerConnect      map[uint64]struct{}
	mxOnPlayerConnect    sync.Mutex
	onPlayerDisconnect   map[uint64]struct{}
	mxOnPlayerDisconnect sync.Mutex
	onPlayerChange       map[uint64]entity.Player
	mxOnPlayerChange     sync.Mutex
}

type config interface {
	NCRequestPerSecond() int
}

func NewProducer(config config, server *netcode.Server) *Producer {
	return &Producer{
		config: config,
		server: server,

		onPlayerConnect:    make(map[uint64]struct{}),
		onPlayerDisconnect: make(map[uint64]struct{}),
		onPlayerChange:     make(map[uint64]entity.Player),
	}
}

func (p *Producer) prepareMessage(code api.Code, payload proto.Message) (*api.Message, error) {
	msg := &api.Message{Code: code}
	var err error
	if msg.Payload, err = proto.Marshal(payload); err != nil {
		return nil, err
	}
	return msg, nil
}

func (p *Producer) marshalMessage(code api.Code, payload proto.Message) ([]byte, error) {
	msg, err := p.prepareMessage(code, payload)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(msg)
}
