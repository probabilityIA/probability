package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueueBoldRequests = rabbitmq.QueuePayBoldRequests

// BoldConsumer consume solicitudes de pago Bold
type BoldConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

// New crea un nuevo consumer Bold
func New(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	logger log.ILogger,
) *BoldConsumer {
	return &BoldConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("bold.consumer"),
	}
}

// Start inicia el consumo de solicitudes Bold
func (c *BoldConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", QueueBoldRequests).Msg("Starting Bold consumer")

	if err := c.rabbit.DeclareQueue(QueueBoldRequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QueueBoldRequests, err)
	}

	if err := c.rabbit.Consume(ctx, QueueBoldRequests, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueueBoldRequests).Msg("Bold consumer started")
	return nil
}

func (c *BoldConsumer) handleMessage(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var msg app.PaymentRequestMsg
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal Bold request")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Bold payment request")

	if err := c.useCase.ProcessPayment(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process Bold payment")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Int64("processing_ms", time.Since(startTime).Milliseconds()).
		Msg("Bold payment processed")

	return nil
}
