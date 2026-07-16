package queue

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type ecommerceStockPushMessage struct {
	ProductID           string `json:"product_id"`
	ExternalProductID   string `json:"external_product_id"`
	IntegrationID       uint   `json:"integration_id"`
	IntegrationTypeCode string `json:"integration_type_code"`
	BusinessID          uint   `json:"business_id"`
	Quantity            int    `json:"quantity"`
	Timestamp           string `json:"timestamp"`
}

type InventoryPushConsumer struct {
	queue   rabbitmq.IQueue
	useCase usecases.IShopifyUseCase
	logger  log.ILogger
}

func NewInventoryPushConsumer(queue rabbitmq.IQueue, useCase usecases.IShopifyUseCase, logger log.ILogger) *InventoryPushConsumer {
	return &InventoryPushConsumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("shopify.inventory_push_consumer"),
	}
}

func (c *InventoryPushConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueShopifyInventoryStockPush, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de push de stock Shopify")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueShopifyInventoryStockPush, func(body []byte) error {
			c.handle(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de push de stock Shopify")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de push de stock Shopify iniciado")
}

func (c *InventoryPushConsumer) handle(ctx context.Context, body []byte) {
	var msg ecommerceStockPushMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Mensaje de push de stock Shopify invalido")
		return
	}

	if msg.ExternalProductID == "" || msg.IntegrationID == 0 {
		c.logger.Warn(ctx).
			Str("product_id", msg.ProductID).
			Uint("integration_id", msg.IntegrationID).
			Msg("Mensaje de push de stock incompleto, se omite")
		return
	}

	integrationID := strconv.FormatUint(uint64(msg.IntegrationID), 10)
	if err := c.useCase.UpdateInventory(ctx, integrationID, msg.ExternalProductID, msg.Quantity); err != nil {
		c.logger.Error(ctx).
			Err(err).
			Str("integration_id", integrationID).
			Str("external_product_id", msg.ExternalProductID).
			Int("quantity", msg.Quantity).
			Msg("Error al empujar stock a Shopify")
	}
}
