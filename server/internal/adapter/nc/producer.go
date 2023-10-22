package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"github.com/dqso/mincer/server/internal/log"
	"github.com/wirepair/netcode"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"sync"
)

type Producer struct {
	config config
	logger log.Logger
	server *netcode.Server

	mxOnPlayerConnect sync.Mutex
	onPlayerConnect   map[uint64]struct{}

	mxOnPlayerDisconnect sync.Mutex
	onPlayerDisconnect   map[uint64]struct{}

	mxOnPlayerWasted sync.Mutex
	onPlayerWasted   []*api.OnPlayerWasted

	mxOnPlayerAttacked sync.Mutex
	onPlayerAttacked   []*api.OnPlayerAttacked

	spawnPlayer   map[uint64]entity.Player
	mxSpawnPlayer sync.Mutex

	mxPlayerStats sync.Mutex
	playerStats   map[uint64]entity.PlayerStats

	mxPlayerHP sync.Mutex
	playerHP   map[uint64]int32

	mxPlayerPositions sync.Mutex
	playerPositions   map[uint64]entity.Point

	mxPlayerWeapon sync.Mutex
	playerWeapon   map[uint64]entity.Weapon

	mxCreateProjectile sync.Mutex
	createProjectile   map[uint64]entity.Projectile

	mxProjectilePositions sync.Mutex
	projectilePositions   map[uint64]entity.Point

	mxDeleteProjectile sync.Mutex
	deleteProjectile   map[uint64]struct{}
}

type config interface {
	NCRequestPerSecond() int
}

func NewProducer(config config, logger log.Logger, server *netcode.Server) *Producer {
	return &Producer{
		config: config,
		logger: logger.With(log.Module("nc_producer")),
		server: server,

		onPlayerConnect:     make(map[uint64]struct{}),
		onPlayerDisconnect:  make(map[uint64]struct{}),
		spawnPlayer:         make(map[uint64]entity.Player),
		playerStats:         make(map[uint64]entity.PlayerStats),
		playerHP:            make(map[uint64]int32),
		playerPositions:     make(map[uint64]entity.Point),
		playerWeapon:        make(map[uint64]entity.Weapon),
		createProjectile:    make(map[uint64]entity.Projectile),
		projectilePositions: make(map[uint64]entity.Point),
		deleteProjectile:    make(map[uint64]struct{}),
	}
}

func (p *Producer) appendToBatch(batch []*api.Message, code api.Code, payload proto.Message) []*api.Message {
	msg := &api.Message{Code: code}
	var err error
	if msg.Payload, err = proto.Marshal(payload); err != nil {
		p.logger.Error("unable to marshal the message", slog.String("code", code.String()), log.Err(err))
		return batch
	}
	return append(batch, msg)
}

// deprecated
func (p *Producer) prepareMessage(code api.Code, payload proto.Message) (*api.Message, error) {
	msg := &api.Message{Code: code}
	var err error
	if msg.Payload, err = proto.Marshal(payload); err != nil {
		return nil, err
	}
	return msg, nil
}

func (p *Producer) marshalMessage(code api.Code, payload proto.Message) ([]byte, error) {
	msg := &api.Message{Code: code}
	var err error
	if msg.Payload, err = proto.Marshal(payload); err != nil {
		return nil, err
	}
	return proto.Marshal(msg)
}
