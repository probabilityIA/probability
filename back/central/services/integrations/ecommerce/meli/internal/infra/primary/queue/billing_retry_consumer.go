package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const maxBillingRetryAttempts = 4

type billingRetryMessage struct {
	IntegrationID string `json:"integration_id"`
	OrderID       int64  `json:"order_id"`
	Attempts      int    `json:"attempts"`
}

type BillingRetryConsumer struct {
	queue  rabbitmq.IQueue
	uc     usecases.IMeliUseCase
	logger log.ILogger
}

func NewBillingRetryConsumer(queue rabbitmq.IQueue, uc usecases.IMeliUseCase, logger log.ILogger) *BillingRetryConsumer {
	return &BillingRetryConsumer{
		queue:  queue,
		uc:     uc,
		logger: logger.WithModule("meli.billing_retry"),
	}
}

func (c *BillingRetryConsumer) Start(ctx context.Context) {
	if c.queue == nil {
		c.logger.Warn(ctx).Msg("RabbitMQ not available, MercadoLibre billing retry consumer disabled")
		return
	}
	if err := c.queue.DeclareQueue(rabbitmq.QueueMeliBillingRetry, true); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to declare billing retry queue")
		return
	}

	c.logger.Info(ctx).Str("queue", rabbitmq.QueueMeliBillingRetry).Msg("Starting MercadoLibre billing retry consumer")

	go func() {
		err := c.queue.Consume(ctx, rabbitmq.QueueMeliBillingRetry, func(body []byte) error {
			c.handleMessage(ctx, body)
			return nil
		})
		if err != nil {
			c.logger.Error(ctx).Err(err).Msg("MercadoLibre billing retry consumer stopped with error")
		}
	}()
}

func (c *BillingRetryConsumer) handleMessage(ctx context.Context, body []byte) {
	var msg billingRetryMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to unmarshal billing retry message")
		return
	}
	if msg.OrderID == 0 || msg.IntegrationID == "" {
		return
	}

	c.wait(ctx, msg.Attempts)

	found, err := c.uc.RetryBilling(ctx, msg.IntegrationID, msg.OrderID)
	if err != nil {
		c.logger.Warn(ctx).Err(err).Int64("order_id", msg.OrderID).Msg("Billing retry attempt failed")
	}
	if found {
		c.logger.Info(ctx).Int64("order_id", msg.OrderID).Msg("Billing info resolved, order republished")
		return
	}

	next := msg.Attempts + 1
	if next > maxBillingRetryAttempts {
		c.logger.Info(ctx).Int64("order_id", msg.OrderID).Int("attempts", msg.Attempts).Msg("Billing info still missing, giving up")
		return
	}
	c.reenqueue(ctx, msg.IntegrationID, msg.OrderID, next)
}

func (c *BillingRetryConsumer) wait(ctx context.Context, attempts int) {
	delay := time.Duration(attempts) * 20 * time.Second
	if delay > 90*time.Second {
		delay = 90 * time.Second
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

func (c *BillingRetryConsumer) reenqueue(ctx context.Context, integrationID string, orderID int64, attempts int) {
	body, err := json.Marshal(billingRetryMessage{
		IntegrationID: integrationID,
		OrderID:       orderID,
		Attempts:      attempts,
	})
	if err != nil {
		return
	}
	if perr := c.queue.Publish(ctx, rabbitmq.QueueMeliBillingRetry, body); perr != nil {
		c.logger.Warn(ctx).Err(perr).Int64("order_id", orderID).Msg("Failed to re-enqueue billing retry")
	}
}
