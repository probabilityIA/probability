package queue

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type productSyncRequest struct {
	IntegrationID uint   `json:"integration_id"`
	BusinessID    uint   `json:"business_id"`
	CorrelationID string `json:"correlation_id"`
}

type ProductSyncConsumer struct {
	queue   rabbitmq.IQueue
	useCase usecases.IWooCommerceUseCase
	logger  log.ILogger
}

func NewProductSyncConsumer(queue rabbitmq.IQueue, useCase usecases.IWooCommerceUseCase, logger log.ILogger) *ProductSyncConsumer {
	return &ProductSyncConsumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("woocommerce"),
	}
}

func (c *ProductSyncConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueWooProductSyncRequests, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de sync de productos WooCommerce")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueWooProductSyncRequests, func(body []byte) error {
			c.handle(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de sync de productos WooCommerce")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de sync de productos WooCommerce iniciado")
}

func (c *ProductSyncConsumer) handle(ctx context.Context, body []byte) {
	var msg productSyncRequest
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Mensaje de sync de productos invalido")
		return
	}

	if msg.IntegrationID == 0 || msg.BusinessID == 0 {
		c.logger.Warn(ctx).Msg("Mensaje de sync de productos incompleto, se omite")
		return
	}

	integrationID := strconv.FormatUint(uint64(msg.IntegrationID), 10)
	if err := c.useCase.SyncProducts(ctx, integrationID, msg.BusinessID, msg.CorrelationID); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("integration_id", integrationID).
			Uint("business_id", msg.BusinessID).
			Msg("Error al sincronizar productos a WooCommerce")
	}
}
