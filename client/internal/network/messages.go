package network

import (
	"github.com/dqso/mincer/client/internal/api"
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/wirepair/netcode"
	"google.golang.org/protobuf/proto"
	"log"
)

func (m *Manager) decodeMessage(data []byte) {
	var msg api.Message
	if err := proto.Unmarshal(data, &msg); err != nil {
		log.Print(err) // TODO logger
		return
	}
	m.decodeMessageWithCode(msg.Code, msg.Payload)
}

func (m *Manager) decodeMessageWithCode(code api.Code, data []byte) {
	switch code {

	case api.Code_BATCH:
		var batch api.Batch
		if err := proto.Unmarshal(data, &batch); err != nil {
			log.Print(err) // TODO logger
			return
		}
		for _, msg := range batch.Messages {
			m.decodeMessageWithCode(msg.Code, msg.Payload)
		}

	case api.Code_ON_PLAYER_CONNECT:
		var message api.OnPlayerConnect
		if err := proto.Unmarshal(data, &message); err != nil {
			log.Print(err) // TODO logger
			return
		}
		log.Printf("Присоединился %d", message.PlayerId)

	case api.Code_ON_PLAYER_DISCONNECT:
		var message api.OnPlayerDisconnect
		if err := proto.Unmarshal(data, &message); err != nil {
			log.Print(err) // TODO logger
			return
		}
		log.Printf("Отсоединился %d", message.PlayerId)
		m.world.Players().Remove(message.PlayerId)

	case api.Code_PLAYER_LIST:
		var message api.PlayerList
		if err := proto.Unmarshal(data, &message); err != nil {
			log.Print(err) // TODO logger
			return
		}
		for _, p := range message.Players {
			createOrChangePlayer(m.world.Players(), p)
		}

	case api.Code_SPAWN_PLAYER:
		var msg api.SpawnPlayer
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Print(err) // TODO logger
			return
		}
		log.Printf("SpawnPlayer: %v", msg.String())
		createOrChangePlayer(m.world.Players(), msg.Player)
	//log.Printf("Изменился %d: (%0.2f, %0.2f), здоровье %d, радиус %0.2f, умер %v", message.PlayerId, message.X, message.Y, message.Hp, message.Radius, message.Dead)

	case api.Code_SET_PLAYER_CLASS:
		var msg api.SetPlayerClass
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Print(err) // TODO logger
			return
		}
		p, ok := m.world.Players().Get(msg.Id)
		if !ok {
			return
		}
		p.SetClass(entity.Class(msg.Class))

	case api.Code_SET_PLAYER_HP:
		var msg api.SetPlayerHP
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Print(err) // TODO logger
			return
		}
		p, ok := m.world.Players().Get(msg.Id)
		if !ok {
			return
		}
		p.SetHP(msg.Hp)

	case api.Code_SET_PLAYER_RADIUS:
		var msg api.SetPlayerRadius
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Print(err) // TODO logger
			return
		}
		p, ok := m.world.Players().Get(msg.Id)
		if !ok {
			return
		}
		p.SetRadius(msg.Radius)

	case api.Code_SET_PLAYER_SPEED:
		var msg api.SetPlayerSpeed
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Print(err) // TODO logger
			return
		}
		p, ok := m.world.Players().Get(msg.Id)
		if !ok {
			return
		}
		p.SetSpeed(msg.Speed)

	case api.Code_SET_PLAYER_POSITION:
		var msg api.SetPlayerPosition
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Print(err) // TODO logger
			return
		}
		p, ok := m.world.Players().Get(msg.Id)
		if !ok {
			return
		}
		p.SetPosition(msg.X, msg.Y)

	default:
		log.Printf("unknown message %s", code.String())
		return
	}
}

func createOrChangePlayer(players entity.Players, p *api.Player) {
	player, ok := players.Get(p.Id)
	if !ok {
		player = entity.NewPlayer(p.Id, entity.Class(p.Class), p.Hp, p.Radius, p.Speed, p.X, p.Y)
		players.Add(player)
	} else {
		player.SetClass(entity.Class(p.Class))
		player.SetHP(p.Hp)
		player.SetRadius(p.Radius)
		player.SetSpeed(p.Speed)
		player.SetPosition(p.X, p.Y)
	}
}

func (m *Manager) repeatingMessageSend() error {
	direction, isMoving := m.world.Players().Me().Direction()

	var err error
	msg := &api.Message{Code: api.Code_CLIENT_INFO}
	msg.Payload, err = proto.Marshal(&api.ClientInfo{
		Direction: direction,
		IsMoving:  isMoving,
	})
	if err != nil {
		return err
	}
	bts, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	if err := m.nc.SendData(bts); err != nil {
		return err
	}
	return nil
}

func (m *Manager) disconnect() error {
	var err error
	msg := &api.Message{Code: api.Code_QUIT}
	msg.Payload, err = proto.Marshal(&api.Quit{})
	if err != nil {
		return err
	}
	bts, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	if err := m.nc.SendData(bts); err != nil {
		return err
	}
	// TODO state to const
	if err := m.nc.Disconnect(netcode.ClientState(100), true); err != nil {
		return err
	}
	return nil
}
