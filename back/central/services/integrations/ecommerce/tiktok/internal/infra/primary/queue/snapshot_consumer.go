package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
	secondaryqueue "github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/infra/secondary/queue"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type SnapshotConsumer struct {
	queue   rabbitmq.IQueue
	useCase usecases.ITikTokUseCase
	logger  log.ILogger
}

func NewSnapshotConsumer(queue rabbitmq.IQueue, useCase usecases.ITikTokUseCase, logger log.ILogger) *SnapshotConsumer {
	return &SnapshotConsumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("tiktok.snapshot_consumer"),
	}
}

func (c *SnapshotConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueTiktokShopSnapshots, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de snapshots de TikTok")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueTiktokShopSnapshots, func(body []byte) error {
			var msg secondaryqueue.SnapshotMessage
			if err := json.Unmarshal(body, &msg); err != nil {
				c.logger.Error(ctx).Err(err).Msg("Mensaje de snapshot de TikTok invalido, descartado")
				return nil
			}
			bg := context.Background()
			snapshot := &domain.ShopSnapshot{
				BusinessID:           msg.BusinessID,
				IntegrationID:        msg.IntegrationID,
				ShopID:               msg.ShopID,
				ShopName:             msg.ShopName,
				ShopRating:           msg.ShopRating,
				ShopReviewCount:      msg.ShopReviewCount,
				ItemSoldCount:        msg.ItemSoldCount,
				ShopPerformanceValue: msg.ShopPerformanceValue,
			}
			if err := c.useCase.SaveSnapshot(bg, snapshot); err != nil {
				c.logger.Error(bg).Err(err).
					Uint("integration_id", msg.IntegrationID).
					Msg("Error procesando snapshot de TikTok")
			}
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de snapshots de TikTok")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de snapshots de TikTok iniciado")
}
