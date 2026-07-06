package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type ProductUpserter interface {
	UpsertFromProvider(ctx context.Context, dto *domain.ProductProviderUpsertDTO) error
}

type productUpsertMessage struct {
	BusinessID     uint    `json:"business_id"`
	IntegrationID  uint    `json:"integration_id"`
	SKU            string  `json:"sku"`
	Name           string  `json:"name"`
	TrackInventory bool    `json:"track_inventory"`
	Price          float64 `json:"price"`
	ExternalID     string  `json:"external_id"`
}

type ProductUpsertConsumer struct {
	queue   rabbitmq.IQueue
	useCase ProductUpserter
	logger  log.ILogger
}

func NewProductUpsertConsumer(queue rabbitmq.IQueue, useCase ProductUpserter, logger log.ILogger) *ProductUpsertConsumer {
	return &ProductUpsertConsumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("products.upsert_consumer"),
	}
}

func (c *ProductUpsertConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueProductsProviderUpsert, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de upsert de productos")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueProductsProviderUpsert, func(body []byte) error {
			c.handle(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de upsert de productos")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de upsert de productos iniciado")
}

func (c *ProductUpsertConsumer) handle(ctx context.Context, body []byte) {
	var msg productUpsertMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Mensaje de upsert de producto invalido")
		return
	}

	if msg.BusinessID == 0 || msg.SKU == "" {
		c.logger.Warn(ctx).Str("sku", msg.SKU).Uint("business_id", msg.BusinessID).Msg("Mensaje de upsert incompleto, se omite")
		return
	}

	err := c.useCase.UpsertFromProvider(ctx, &domain.ProductProviderUpsertDTO{
		BusinessID:     msg.BusinessID,
		SKU:            msg.SKU,
		Name:           msg.Name,
		TrackInventory: msg.TrackInventory,
		Price:          msg.Price,
		ExternalID:     msg.ExternalID,
	})
	if err != nil {
		c.logger.Error(ctx).Err(err).Str("sku", msg.SKU).Uint("business_id", msg.BusinessID).Msg("Error al hacer upsert de producto desde proveedor")
	}
}
