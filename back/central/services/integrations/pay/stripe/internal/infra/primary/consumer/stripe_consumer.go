package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/stripe/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueueStripeRequests = "pay.stripe.requests"

// StripeConsumer consume solicitudes de pago Stripe
type StripeConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

// New crea un nuevo consumer Stripe
func New(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	logger log.ILogger,
) *StripeConsumer {
	return &StripeConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("stripe.consumer"),
	}
}

// Start inicia el consumo de solicitudes Stripe
func (c *StripeConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", QueueStripeRequests).Msg("Starting Stripe consumer")

	if err := c.rabbit.DeclareQueue(QueueStripeRequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QueueStripeRequests, err)
	}

	if err := c.rabbit.Consume(ctx, QueueStripeRequests, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueueStripeRequests).Msg("Stripe consumer started")
	return nil
}

func (c *StripeConsumer) handleMessage(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var msg app.PaymentRequestMsg
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal Stripe request")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Stripe payment request")

	if err := c.useCase.ProcessPayment(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process Stripe payment")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Int64("processing_ms", time.Since(startTime).Milliseconds()).
		Msg("Stripe payment processed")

	return nil
}
