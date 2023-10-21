package nc_adapter

import (
	"context"
	"github.com/dqso/mincer/server/internal/api"
	"log"
	"time"
)

func (p *Producer) StartLoop(ctx context.Context) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second / time.Duration(p.config.NCRequestPerSecond())):
			}

			batch := new(api.Batch)
			batch.Messages = append(batch.Messages, p.onPlayerConnectBatch()...)
			batch.Messages = append(batch.Messages, p.onPlayerDisconnectBatch()...)
			batch.Messages = append(batch.Messages, p.onPlayerWastedBatch()...)
			batch.Messages = append(batch.Messages, p.spawnPlayerBatch()...)
			batch.Messages = append(batch.Messages, p.setPlayerStatsBatch()...)
			batch.Messages = append(batch.Messages, p.setPlayerHPBatch()...)
			batch.Messages = append(batch.Messages, p.setPlayerWeaponBatch()...)
			batch.Messages = append(batch.Messages, p.setPlayerPositionBatch()...)
			batch.Messages = append(batch.Messages, p.createProjectileBatch()...)
			batch.Messages = append(batch.Messages, p.setProjectilePositionBatch()...)
			batch.Messages = append(batch.Messages, p.deleteProjectileBatch()...)

			bts, err := p.marshalMessage(api.Code_BATCH, batch)
			if err != nil {
				log.Print(err) // TODO logger
				continue
			}
			p.server.SendPayloads(bts)
		}
	}()
	return done
}
