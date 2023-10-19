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

	mxOnPlayerConnect sync.Mutex
	onPlayerConnect   map[uint64]struct{}

	mxOnPlayerDisconnect sync.Mutex
	onPlayerDisconnect   map[uint64]struct{}

	spawnPlayer   map[uint64]entity.Player
	mxSpawnPlayer sync.Mutex

	mxPlayerClasses sync.Mutex
	playerClasses   map[uint64]entity.Class

	mxPlayerHP sync.Mutex
	playerHP   map[uint64]int64

	mxPlayerRadius sync.Mutex
	playerRadius   map[uint64]float64

	mxPlayerSpeed sync.Mutex
	playerSpeed   map[uint64]float64

	mxPlayerPositions sync.Mutex
	playerPositions   map[uint64]entity.Point
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
		spawnPlayer:        make(map[uint64]entity.Player),
		playerClasses:      make(map[uint64]entity.Class),
		playerHP:           make(map[uint64]int64),
		playerRadius:       make(map[uint64]float64),
		playerSpeed:        make(map[uint64]float64),
		playerPositions:    make(map[uint64]entity.Point),
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
