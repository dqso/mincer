package network

import (
	"github.com/dqso/mincer/client/internal/api"
	"github.com/dqso/mincer/client/internal/entity"
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

	case api.Code_ON_PLAYER_CHANGE:
		var msg api.OnPlayerChange
		if err := proto.Unmarshal(data, &msg); err != nil {
			log.Print(err) // TODO logger
			return
		}
		p, ok := m.world.Players().Get(msg.Player.PlayerId)
		if !ok {
			p = entity.NewPlayer(msg.Player.PlayerId, msg.Player.X, msg.Player.Y, msg.Player.Hp, msg.Player.Radius, msg.Player.Dead)
			m.world.Players().Add(p)
		} else {
			p.SetPosition(msg.Player.X, msg.Player.Y)
			p.SetStats(msg.Player.Hp, msg.Player.Radius, msg.Player.Dead)
		}
	//log.Printf("Изменился %d: (%0.2f, %0.2f), здоровье %d, радиус %0.2f, умер %v", message.PlayerId, message.X, message.Y, message.Hp, message.Radius, message.Dead)

	case api.Code_PLAYER_LIST:
		var message api.PlayerList
		if err := proto.Unmarshal(data, &message); err != nil {
			log.Print(err) // TODO logger
			return
		}
		for _, player := range message.Players {
			p, ok := m.world.Players().Get(player.PlayerId)
			if !ok {
				p = entity.NewPlayer(player.PlayerId, player.X, player.Y, player.Hp, player.Radius, player.Dead)
				m.world.Players().Add(p)
			} else {
				p.SetPosition(player.X, player.Y)
				p.SetStats(player.Hp, player.Radius, player.Dead)
			}
		}

	default:
		log.Printf("unknown message %s", code.String())
		return
	}
}
