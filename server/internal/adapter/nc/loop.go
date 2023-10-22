package nc_adapter

import (
	"context"
	"github.com/dqso/mincer/server/internal/api"
	"github.com/dqso/mincer/server/internal/log"
	"log/slog"
	"time"
)

func (p *Producer) StartLoop(ctx context.Context) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		p.logger.Debug("started sending periodic messages")
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
			batch.Messages = append(batch.Messages, p.onPlayerAttackedBatch()...)
			batch.Messages = append(batch.Messages, p.spawnPlayerBatch()...)
			batch.Messages = append(batch.Messages, p.setPlayerStatsBatch()...)
			batch.Messages = append(batch.Messages, p.setPlayerHPBatch()...)
			batch.Messages = append(batch.Messages, p.setPlayerWeaponBatch()...)
			batch.Messages = append(batch.Messages, p.setPlayerPositionBatch()...)
			batch.Messages = append(batch.Messages, p.createProjectileBatch()...)
			batch.Messages = append(batch.Messages, p.setProjectilePositionBatch()...)
			batch.Messages = append(batch.Messages, p.deleteProjectileBatch()...)

			const code = api.Code_BATCH
			bts, err := p.marshalMessage(code, batch)
			if err != nil {
				p.logger.Error("unable to marshal the message", slog.String("code", code.String()), log.Err(err))
				continue
			}
			p.server.SendPayloads(bts)
			if len(batch.Messages) > 0 {
				p.logger.Debug("repeating message has been sent to everyone",
					slog.Int("size_message", len(bts)),
				)
			}
		}
	}()
	return done
}
