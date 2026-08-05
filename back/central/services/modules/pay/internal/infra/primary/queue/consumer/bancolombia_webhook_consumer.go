package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type BancolombiaWebhookConsumer struct {
	queue   rabbitmq.IQueue
	useCase ports.IUseCase
	log     log.ILogger
}

func NewBancolombiaWebhookConsumer(queue rabbitmq.IQueue, useCase ports.IUseCase, logger log.ILogger) *BancolombiaWebhookConsumer {
	return &BancolombiaWebhookConsumer{
		queue:   queue,
		useCase: useCase,
		log:     logger.WithModule("pay.bancolombia_webhook_consumer"),
	}
}

func (c *BancolombiaWebhookConsumer) Start(ctx context.Context) error {
	if c.queue == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", rabbitmq.QueuePayBancolombiaWebhookEvents).
		Msg("Starting Bancolombia webhook consumer")

	if err := c.queue.DeclareQueue(rabbitmq.QueuePayBancolombiaWebhookEvents, true); err != nil {
		return fmt.Errorf("declare %s: %w", rabbitmq.QueuePayBancolombiaWebhookEvents, err)
	}

	return c.queue.Consume(ctx, rabbitmq.QueuePayBancolombiaWebhookEvents, c.handleMessage)
}

func (c *BancolombiaWebhookConsumer) handleMessage(body []byte) error {
	ctx := context.Background()

	var msg dtos.BancolombiaWebhookMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.log.Error(ctx).Err(err).Msg("bancolombia webhook consumer: invalid message")
		return nil
	}

	if err := c.useCase.ProcessBancolombiaWebhookMessage(ctx, &msg); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("event_id", msg.EventID).
			Msg("bancolombia webhook consumer: processing failed")
		return err
	}

	return nil
}
