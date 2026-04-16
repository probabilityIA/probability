package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/payu/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueuePayURequests = rabbitmq.QueuePayPayURequests

// PayUConsumer consume solicitudes de pago PayU
type PayUConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

// New crea un nuevo consumer PayU
func New(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	logger log.ILogger,
) *PayUConsumer {
	return &PayUConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("payu.consumer"),
	}
}

// Start inicia el consumo de solicitudes PayU
func (c *PayUConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", QueuePayURequests).Msg("Starting PayU consumer")

	if err := c.rabbit.DeclareQueue(QueuePayURequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QueuePayURequests, err)
	}

	if err := c.rabbit.Consume(ctx, QueuePayURequests, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueuePayURequests).Msg("PayU consumer started")
	return nil
}

func (c *PayUConsumer) handleMessage(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var msg app.PaymentRequestMsg
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal PayU request")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing PayU payment request")

	if err := c.useCase.ProcessPayment(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process PayU payment")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Int64("processing_ms", time.Since(startTime).Milliseconds()).
		Msg("PayU payment processed")

	return nil
}
