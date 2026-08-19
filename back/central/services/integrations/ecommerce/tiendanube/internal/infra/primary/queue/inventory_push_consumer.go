package queue

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
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
	useCase usecases.ITiendanubeUseCase
	logger  log.ILogger
}

func NewInventoryPushConsumer(queue rabbitmq.IQueue, useCase usecases.ITiendanubeUseCase, logger log.ILogger) *InventoryPushConsumer {
	return &InventoryPushConsumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("tiendanube"),
	}
}

func (c *InventoryPushConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueTiendanubeInventoryStockPush, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de push de stock Tiendanube")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueTiendanubeInventoryStockPush, func(body []byte) error {
			return c.handle(ctx, body)
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de push de stock Tiendanube")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de push de stock Tiendanube iniciado")
}

func isPermanent(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	permanentes := []string{
		"integration not found",
		"invalid credentials",
		"missing access_token",
		"missing store_id",
		"no incluye variante",
		"external_id invalido",
		"base_url",
		"returned 404",
		"returned 422",
	}
	for _, frase := range permanentes {
		if strings.Contains(message, frase) {
			return true
		}
	}
	return false
}

func (c *InventoryPushConsumer) handle(ctx context.Context, body []byte) error {
	var msg ecommerceStockPushMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Warn(ctx).Err(err).Msg("Mensaje de push de stock Tiendanube invalido, se descarta")
		return nil
	}

	if msg.ExternalProductID == "" || msg.IntegrationID == 0 {
		c.logger.Warn(ctx).
			Str("product_id", msg.ProductID).
			Uint("integration_id", msg.IntegrationID).
			Msg("Mensaje de push de stock incompleto, se descarta")
		return nil
	}

	integrationID := strconv.FormatUint(uint64(msg.IntegrationID), 10)
	err := c.useCase.UpdateInventory(ctx, integrationID, msg.ExternalProductID, msg.Quantity)
	if err == nil {
		return nil
	}

	if isPermanent(err) {
		c.logger.Warn(ctx).Err(err).
			Str("integration_id", integrationID).
			Str("external_product_id", msg.ExternalProductID).
			Msg("Push de stock a Tiendanube descartado por error permanente (ACK)")
		return nil
	}

	c.logger.Error(ctx).Err(err).
		Str("integration_id", integrationID).
		Str("external_product_id", msg.ExternalProductID).
		Int("quantity", msg.Quantity).
		Msg("Error al empujar stock a Tiendanube, se reintentara")
	return err
}
