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

	case api.Code_ON_PLAYER_WASTED:
		var message api.OnPlayerWasted
		if err := proto.Unmarshal(data, &message); err != nil {
			log.Print(err) // TODO logger
			return
		}
		log.Printf("Игрок %d убит игроком %d", message.Id, message.Killer)
		m.world.AddNewKill(message.Id, message.Killer)

	case api.Code_WORLD_INFO:
		var message api.WorldInfo
		if err := proto.Unmarshal(data, &message); err != nil {
			log.Print(err) // TODO logger
			return
		}
		m.world.SetSize(message.Northwest.X, message.Northwest.Y, message.Southeast.X, message.Southeast.Y)

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

	case api.Code_SET_PLAYER_STATS:
		var msg api.SetPlayerStats
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Print(err) // TODO logger
			return
		}
		p, ok := m.world.Players().Get(msg.Id)
		if !ok {
			return
		}
		p.SetStats(dtoPlayerStats(msg.Stats))

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
		if msg.Hp == 0 && msg.Id == m.world.Players().Me().ID() {
			m.beReborn = make(chan struct{})
		}

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
		p.SetPosition(entity.Point{X: msg.X, Y: msg.Y})

	default:
		log.Printf("unknown message %s", code.String())
		return
	}
}

func createOrChangePlayer(players entity.Players, p *api.Player) {
	stats := dtoPlayerStats(p.Stats)
	player, ok := players.Get(p.Id)
	if !ok {
		player = entity.NewPlayer(p.Id, p.Hp, stats, entity.Point{X: p.X, Y: p.Y})
		players.Add(player)
	} else {
		player.SetStats(stats)
		player.SetHP(p.Hp)
		player.SetPosition(entity.Point{X: p.X, Y: p.Y})
	}
}

func dtoPlayerStats(stats *api.PlayerStats) entity.PlayerStats {
	return entity.NewPlayerStats(
		entity.Class(stats.Class),
		stats.Radius,
		stats.Speed,
		stats.MaxHP,
	)
}

func (m *Manager) repeatingMessageSend() error {
	me := m.world.Players().Me()
	direction, isMoving := me.Direction()

	var err error
	msg := &api.Message{Code: api.Code_CLIENT_INFO}
	msg.Payload, err = proto.Marshal(&api.ClientInfo{
		Direction: direction,
		IsMoving:  isMoving,
		Attack:    me.Attack(),
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

func (m *Manager) sendBeReborn() error {
	var err error
	msg := &api.Message{Code: api.Code_BE_REBORN}
	msg.Payload, err = proto.Marshal(&api.BeReborn{})
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
