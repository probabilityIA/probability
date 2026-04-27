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

type BoldWebhookConsumer struct {
	queue   rabbitmq.IQueue
	useCase ports.IUseCase
	log     log.ILogger
}

func NewBoldWebhookConsumer(queue rabbitmq.IQueue, useCase ports.IUseCase, logger log.ILogger) *BoldWebhookConsumer {
	return &BoldWebhookConsumer{
		queue:   queue,
		useCase: useCase,
		log:     logger.WithModule("pay.bold_webhook_consumer"),
	}
}

func (c *BoldWebhookConsumer) Start(ctx context.Context) error {
	if c.queue == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).
		Str("queue", rabbitmq.QueuePayBoldWebhookEvents).
		Msg("Starting Bold webhook consumer")

	if err := c.queue.DeclareQueue(rabbitmq.QueuePayBoldWebhookEvents, true); err != nil {
		return fmt.Errorf("declare %s: %w", rabbitmq.QueuePayBoldWebhookEvents, err)
	}

	return c.queue.Consume(ctx, rabbitmq.QueuePayBoldWebhookEvents, c.handleMessage)
}

func (c *BoldWebhookConsumer) handleMessage(body []byte) error {
	ctx := context.Background()

	var msg dtos.BoldWebhookMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.log.Error(ctx).Err(err).Msg("bold webhook consumer: invalid message")
		return nil
	}

	if err := c.useCase.ProcessBoldWebhookMessage(ctx, &msg); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("bold_event_id", msg.BoldEventID).
			Msg("bold webhook consumer: processing failed")
		return err
	}

	return nil
}
