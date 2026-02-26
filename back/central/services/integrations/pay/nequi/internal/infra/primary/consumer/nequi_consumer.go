package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueueNequiRequests = "pay.nequi.requests"

// NequiConsumer consume solicitudes de pago Nequi
type NequiConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

// New crea un nuevo consumer Nequi
func New(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	logger log.ILogger,
) *NequiConsumer {
	return &NequiConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("nequi.consumer"),
	}
}

// Start inicia el consumo de solicitudes Nequi
func (c *NequiConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", QueueNequiRequests).Msg("Starting Nequi consumer")

	if err := c.rabbit.DeclareQueue(QueueNequiRequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QueueNequiRequests, err)
	}

	if err := c.rabbit.Consume(ctx, QueueNequiRequests, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueueNequiRequests).Msg("Nequi consumer started")
	return nil
}

func (c *NequiConsumer) handleMessage(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var msg app.PaymentRequestMsg
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal Nequi request")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing Nequi payment request")

	if err := c.useCase.ProcessPayment(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process Nequi payment")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Int64("processing_ms", time.Since(startTime).Milliseconds()).
		Msg("Nequi payment processed")

	return nil
}
