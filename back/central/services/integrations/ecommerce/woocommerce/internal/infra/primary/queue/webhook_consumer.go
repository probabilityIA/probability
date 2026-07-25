package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type WebhookMessage struct {
	Topic         string          `json:"topic"`
	Source        string          `json:"source"`
	IntegrationID string          `json:"integration_id"`
	Body          json.RawMessage `json:"body"`
}

type WebhookConsumer struct {
	queue   rabbitmq.IQueue
	useCase usecases.IWooCommerceUseCase
	logger  log.ILogger
}

func NewWebhookConsumer(queue rabbitmq.IQueue, useCase usecases.IWooCommerceUseCase, logger log.ILogger) *WebhookConsumer {
	return &WebhookConsumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("woocommerce.webhook_consumer"),
	}
}

func (c *WebhookConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueWebhooksWoocommerceReceived, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de webhooks WooCommerce")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueWebhooksWoocommerceReceived, func(body []byte) error {
			var msg WebhookMessage
			if err := json.Unmarshal(body, &msg); err != nil {
				c.logger.Error(ctx).Err(err).Msg("Mensaje de webhook WooCommerce invalido, descartado")
				return nil
			}
			bg := context.Background()
			if err := c.useCase.ProcessWebhookOrder(bg, msg.Topic, msg.Source, msg.IntegrationID, msg.Body); err != nil {
				c.logger.Error(bg).Err(err).
					Str("topic", msg.Topic).
					Str("source", msg.Source).
					Str("integration_id", msg.IntegrationID).
					Msg("Failed to process WooCommerce webhook order")
			}
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de webhooks WooCommerce")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de webhooks WooCommerce iniciado")
}
