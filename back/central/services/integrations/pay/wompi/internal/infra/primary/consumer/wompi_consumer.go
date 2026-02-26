package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/wompi/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueueWompiRequests = "pay.wompi.requests"

// WompiConsumer consume solicitudes de pago Wompi
type WompiConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

// New crea un nuevo consumer Wompi
func New(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	logger log.ILogger,
) *WompiConsumer {
	return &WompiConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("wompi.consumer"),
	}
}

// Start inicia el consumo de solicitudes Wompi
func (c *WompiConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", QueueWompiRequests).Msg("Starting Wompi consumer")

	if err := c.rabbit.DeclareQueue(QueueWompiRequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QueueWompiRequests, err)
	}

	if err := c.rabbit.Consume(ctx, QueueWompiRequests, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueueWompiRequests).Msg("Wompi consumer started")
	return nil
}

func (c *WompiConsumer) handleMessage(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var msg app.PaymentRequestMsg
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal Wompi request")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Wompi payment request")

	if err := c.useCase.ProcessPayment(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process Wompi payment")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Int64("processing_ms", time.Since(startTime).Milliseconds()).
		Msg("Wompi payment processed")

	return nil
}
