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

	case api.Code_ON_PLAYER_CHANGE:
		var message api.OnPlayerChange
		if err := proto.Unmarshal(data, &message); err != nil {
			log.Print(err) // TODO logger
			return
		}
		p, ok := m.world.Player(message.Player.PlayerId)
		if !ok {
			p = &entity.Player{}
		}
		p.ID = message.Player.PlayerId
		p.X = message.Player.X
		p.Y = message.Player.Y
		p.HP = message.Player.Hp
		p.Radius = message.Player.Radius
		p.Dead = message.Player.Dead
		m.world.SetPlayer(p)
	//log.Printf("Изменился %d: (%0.2f, %0.2f), здоровье %d, радиус %0.2f, умер %v", message.PlayerId, message.X, message.Y, message.Hp, message.Radius, message.Dead)

	case api.Code_PLAYER_LIST:
		var message api.PlayerList
		if err := proto.Unmarshal(data, &message); err != nil {
			log.Print(err) // TODO logger
			return
		}
		for _, player := range message.Players {
			p, ok := m.world.Player(player.PlayerId)
			if !ok {
				p = &entity.Player{}
			}
			p.ID = player.PlayerId
			p.X = player.X
			p.Y = player.Y
			p.HP = player.Hp
			p.Radius = player.Radius
			p.Dead = player.Dead
			m.world.SetPlayer(p)
		}

	default:
		log.Printf("unknown message %s", code.String())
		return
	}
}
