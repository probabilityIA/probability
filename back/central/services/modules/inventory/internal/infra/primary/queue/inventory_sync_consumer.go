package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type InventorySyncConsumer struct {
	queue  rabbitmq.IQueue
	uc     app.IUseCase
	logger log.ILogger
}

func NewInventorySyncConsumer(queue rabbitmq.IQueue, uc app.IUseCase, logger log.ILogger) *InventorySyncConsumer {
	return &InventorySyncConsumer{
		queue:  queue,
		uc:     uc,
		logger: logger.WithModule("inventory.provider_sync_consumer"),
	}
}

func (c *InventorySyncConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, provider sync consumer disabled")
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueInventoryProviderSync, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare provider sync queue")
		return
	}

	c.logger.Info(ctx).Str("queue", rabbitmq.QueueInventoryProviderSync).Msg("Starting provider sync consumer")

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueInventoryProviderSync, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Provider sync consumer stopped with error")
		}
	}()
}

func (c *InventorySyncConsumer) handleMessage(ctx context.Context, body []byte) {
	var dto request.ProviderStockSyncDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to unmarshal provider sync message")
		return
	}

	c.logger.Info(ctx).
		Uint("business_id", dto.BusinessID).
		Str("provider", dto.Provider).
		Int("items", len(dto.Items)).
		Msg("Processing provider stock sync")

	if _, err := c.uc.SyncProviderStock(ctx, dto); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Provider stock sync failed")
	}
}
