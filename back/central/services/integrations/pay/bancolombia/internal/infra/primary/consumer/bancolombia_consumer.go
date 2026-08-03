package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueueBancolombiaRequests = rabbitmq.QueuePayBancolombiaRequests

type BancolombiaConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

func New(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	logger log.ILogger,
) *BancolombiaConsumer {
	return &BancolombiaConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("bancolombia.consumer"),
	}
}

func (c *BancolombiaConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", QueueBancolombiaRequests).Msg("Starting Bancolombia consumer")

	if err := c.rabbit.DeclareQueue(QueueBancolombiaRequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QueueBancolombiaRequests, err)
	}

	if err := c.rabbit.Consume(ctx, QueueBancolombiaRequests, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueueBancolombiaRequests).Msg("Bancolombia consumer started")
	return nil
}

func (c *BancolombiaConsumer) handleMessage(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var msg app.PaymentRequestMsg
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal Bancolombia request")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Bancolombia payment request")

	if err := c.useCase.ProcessPayment(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process Bancolombia payment")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Int64("processing_ms", time.Since(startTime).Milliseconds()).
		Msg("Bancolombia payment processed")

	return nil
}
