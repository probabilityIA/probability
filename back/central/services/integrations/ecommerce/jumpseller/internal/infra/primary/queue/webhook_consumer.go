package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type WebhookMessage struct {
	Event         string          `json:"event"`
	StoreCode     string          `json:"store_code"`
	IntegrationID string          `json:"integration_id"`
	Body          json.RawMessage `json:"body"`
}

type WebhookConsumer struct {
	queue   rabbitmq.IQueue
	useCase usecases.IJumpsellerUseCase
	logger  log.ILogger
}

func NewWebhookConsumer(queue rabbitmq.IQueue, useCase usecases.IJumpsellerUseCase, logger log.ILogger) *WebhookConsumer {
	return &WebhookConsumer{
		queue:   queue,
		useCase: useCase,
		logger:  logger.WithModule("jumpseller.webhook_consumer"),
	}
}

func (c *WebhookConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		return
	}

	if err := c.queue.DeclareQueue(rabbitmq.QueueWebhooksJumpsellerReceived, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Error al declarar la cola de webhooks Jumpseller")
		return
	}

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueWebhooksJumpsellerReceived, func(body []byte) error {
			var msg WebhookMessage
			if err := json.Unmarshal(body, &msg); err != nil {
				c.logger.Error(ctx).Err(err).Msg("Mensaje de webhook Jumpseller invalido, descartado")
				return nil
			}
			bg := context.Background()
			if err := c.useCase.ProcessWebhookOrder(bg, msg.Event, msg.StoreCode, msg.IntegrationID, msg.Body); err != nil {
				c.logger.Error(bg).Err(err).
					Str("event", msg.Event).
					Str("store_code", msg.StoreCode).
					Str("integration_id", msg.IntegrationID).
					Msg("Failed to process Jumpseller webhook order")
			}
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("Error al consumir la cola de webhooks Jumpseller")
		}
	}()

	c.logger.Info(ctx).Msg("Consumer de webhooks Jumpseller iniciado")
}
