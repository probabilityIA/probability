package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/epayco/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueueEPaycoRequests = "pay.epayco.requests"

// EPaycoConsumer consume solicitudes de pago ePayco
type EPaycoConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

// New crea un nuevo consumer ePayco
func New(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	logger log.ILogger,
) *EPaycoConsumer {
	return &EPaycoConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("epayco.consumer"),
	}
}

// Start inicia el consumo de solicitudes ePayco
func (c *EPaycoConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", QueueEPaycoRequests).Msg("Starting ePayco consumer")

	if err := c.rabbit.DeclareQueue(QueueEPaycoRequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QueueEPaycoRequests, err)
	}

	if err := c.rabbit.Consume(ctx, QueueEPaycoRequests, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueueEPaycoRequests).Msg("ePayco consumer started")
	return nil
}

func (c *EPaycoConsumer) handleMessage(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var msg app.PaymentRequestMsg
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal ePayco request")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing ePayco payment request")

	if err := c.useCase.ProcessPayment(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process ePayco payment")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Int64("processing_ms", time.Since(startTime).Milliseconds()).
		Msg("ePayco payment processed")

	return nil
}
