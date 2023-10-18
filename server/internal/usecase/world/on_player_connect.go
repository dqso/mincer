package usecase_world

import "log"

func (uc *Usecase) OnPlayerConnect(connect chan uint64, disconnect chan uint64) {
	for {
		select {

		case id, ok := <-connect:
			if !ok {
				return
			}
			player, err := uc.world.AddPlayer(id)
			if err != nil {
				log.Print(err) // TODO logger
				continue
			}
			uc.ncProducer.OnPlayerConnect(player.ID())
			uc.ncProducer.DirectPlayerList(id, uc.world.Players().Slice())
			uc.ncProducer.OnPlayerChange(player)

		case id, ok := <-disconnect:
			if !ok {
				return
			}
			uc.world.RemovePlayer(id)
			uc.ncProducer.OnPlayerDisconnect(id)
		}
	}
}
