package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/melipago/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueueMeliPagoRequests = rabbitmq.QueuePayMeliPagoRequests

// MeliPagoConsumer consume solicitudes de pago MercadoPago
type MeliPagoConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase app.IUseCase
	log     log.ILogger
}

// New crea un nuevo consumer MercadoPago
func New(
	rabbit rabbitmq.IQueue,
	useCase app.IUseCase,
	logger log.ILogger,
) *MeliPagoConsumer {
	return &MeliPagoConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("melipago.consumer"),
	}
}

// Start inicia el consumo de solicitudes MercadoPago
func (c *MeliPagoConsumer) Start(ctx context.Context) error {
	if c.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	c.log.Info(ctx).Str("queue", QueueMeliPagoRequests).Msg("Starting MercadoPago consumer")

	if err := c.rabbit.DeclareQueue(QueueMeliPagoRequests, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QueueMeliPagoRequests, err)
	}

	if err := c.rabbit.Consume(ctx, QueueMeliPagoRequests, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueueMeliPagoRequests).Msg("MercadoPago consumer started")
	return nil
}

func (c *MeliPagoConsumer) handleMessage(message []byte) error {
	ctx := context.Background()
	startTime := time.Now()

	var msg app.PaymentRequestMsg
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal MercadoPago request")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Float64("amount", msg.Amount).
		Str("reference", msg.Reference).
		Msg("Processing MercadoPago payment request")

	if err := c.useCase.ProcessPayment(ctx, &msg); err != nil {
		c.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to process MercadoPago payment")
		return err
	}

	c.log.Info(ctx).
		Uint("transaction_id", msg.PaymentTransactionID).
		Int64("processing_ms", time.Since(startTime).Milliseconds()).
		Msg("MercadoPago payment processed")

	return nil
}
