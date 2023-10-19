package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
	"log"
)

func (p *Producer) PlayerList(toPlayerID uint64, players []entity.Player) {
	msg := &api.PlayerList{
		Players: make([]*api.Player, 0, len(players)),
	}
	for _, player := range players {
		msg.Players = append(msg.Players, dtoPlayerToApiPlayer(player))
	}
	bts, err := p.marshalMessage(api.Code_PLAYER_LIST, msg)
	if err != nil {
		log.Print(err) // TODO logger
		return
	}
	if err := p.SendPayloadToClient(toPlayerID, bts); err != nil {
		log.Print(err) // TODO logger
		return
	}
}
