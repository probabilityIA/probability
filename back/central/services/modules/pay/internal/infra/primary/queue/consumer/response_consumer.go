package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/constants"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// ResponseConsumer consume respuestas del gateway desde pay.responses
type ResponseConsumer struct {
	queue   rabbitmq.IQueue
	useCase ports.IUseCase
	log     log.ILogger
}

// NewResponseConsumer crea un nuevo consumer de respuestas de pagos
func NewResponseConsumer(
	queue rabbitmq.IQueue,
	useCase ports.IUseCase,
	logger log.ILogger,
) *ResponseConsumer {
	return &ResponseConsumer{
		queue:   queue,
		useCase: useCase,
		log:     logger.WithModule("pay.response_consumer"),
	}
}

// Start inicia el consumo de respuestas de pago
func (c *ResponseConsumer) Start(ctx context.Context) error {
	if c.queue == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", constants.QueuePayResponses).Msg("Starting pay response consumer")

	if err := c.queue.DeclareQueue(constants.QueuePayResponses, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", constants.QueuePayResponses, err)
	}

	if err := c.queue.Consume(ctx, constants.QueuePayResponses, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", constants.QueuePayResponses).Msg("Pay response consumer started")
	return nil
}

func (c *ResponseConsumer) handleMessage(message []byte) error {
	ctx := context.Background()

	var msg dtos.PaymentResponseMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal payment response")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("status", msg.Status).
		Str("gateway", msg.GatewayCode).
		Msg("Received payment response")

	if err := c.useCase.ProcessPaymentResponse(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process payment response")
		return err
	}

	return nil
}
