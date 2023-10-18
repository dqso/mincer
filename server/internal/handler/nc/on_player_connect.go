package nc_handler

import (
	"context"
)

func onPlayerConnectDetector(ctx context.Context, processedPlayersID chan []uint64) (chan uint64, chan uint64) {
	connect, disconnect := make(chan uint64, 1), make(chan uint64, 1)

	go func() {
		defer close(connect)
		defer close(disconnect)
		players := make(map[uint64]int)
		for {
			select {
			case <-ctx.Done():
				return

			case ids, ok := <-processedPlayersID:
				if !ok {
					return
				}
				for _, id := range ids {
					if _, ok := players[id]; !ok {
						connect <- id
					}
					players[id] = 10
				}
				for id, priority := range players {
					if priority <= 0 {
						delete(players, id)
						disconnect <- id
						continue
					}
					players[id]--
				}
			}
		}
	}()

	return connect, disconnect
}
