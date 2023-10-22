package nc_adapter

import (
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/entity"
)

func (p *Producer) PlayerList(toPlayerID uint64, players []entity.Player) {
	const N = 5
	done := make(chan struct{})
	chPlayers := make(chan *api.Player, 1)
	go func() {
		defer close(done)
		for {
			msg := &api.PlayerList{
				Players: make([]*api.Player, 0, N),
			}
			for player := range chPlayers {
				msg.Players = append(msg.Players, player)
				if len(msg.Players) >= N {
					break
				}
			}
			if len(msg.Players) == 0 {
				return
			}
			if !p.SendPayloadToClient(toPlayerID, api.Code_PLAYER_LIST, msg) {
				return
			}
		}
	}()
	for _, player := range players {
		chPlayers <- dtoPlayerToApiPlayer(player)
	}
	close(chPlayers)
	<-done
}
