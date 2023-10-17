package nc_handler

import "context"

func onPlayerConnectDetector(ctx context.Context, processedPlayersID chan uint64) chan uint64 {
	newPlayers := make(chan uint64, 1)
	players := make(map[uint64]struct{})

	go func() {
		defer close(newPlayers)
		for {
			select {
			case <-ctx.Done():
				return
			case id, ok := <-processedPlayersID:
				if !ok {
					return
				}
				if _, ok := players[id]; ok {
					continue
				}
				players[id] = struct{}{}
				newPlayers <- id
			}
		}
	}()

	return newPlayers
}
