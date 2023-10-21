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

	mxOnPlayerWasted sync.Mutex
	onPlayerWasted   []*api.OnPlayerWasted

	spawnPlayer   map[uint64]entity.Player
	mxSpawnPlayer sync.Mutex

	mxPlayerStats sync.Mutex
	playerStats   map[uint64]*api.PlayerStats

	mxPlayerHP sync.Mutex
	playerHP   map[uint64]int32

	mxPlayerWeapon sync.Mutex
	playerWeapon   map[uint64]*api.Weapon

	mxCreateProjectile sync.Mutex
	createProjectile   map[uint64]*api.CreateProjectile

	mxPlayerPositions sync.Mutex
	playerPositions   map[uint64]entity.Point

	mxProjectilePositions sync.Mutex
	projectilePositions   map[uint64]entity.Point

	mxDeleteProjectile sync.Mutex
	deleteProjectile   map[uint64]struct{}
}

type config interface {
	NCRequestPerSecond() int
}

func NewProducer(config config, server *netcode.Server) *Producer {
	return &Producer{
		config: config,
		server: server,

		onPlayerConnect:     make(map[uint64]struct{}),
		onPlayerDisconnect:  make(map[uint64]struct{}),
		spawnPlayer:         make(map[uint64]entity.Player),
		playerStats:         make(map[uint64]*api.PlayerStats),
		playerHP:            make(map[uint64]int32),
		playerWeapon:        make(map[uint64]*api.Weapon),
		createProjectile:    make(map[uint64]*api.CreateProjectile),
		playerPositions:     make(map[uint64]entity.Point),
		projectilePositions: make(map[uint64]entity.Point),
		deleteProjectile:    make(map[uint64]struct{}),
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
