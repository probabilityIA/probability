package router

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	QueuePayRequests   = "pay.requests"
	QueueNequiRequests = "pay.nequi.requests"
)

// payRequestHeader contiene solo los campos para enrutar
type payRequestHeader struct {
	PaymentTransactionID uint      `json:"payment_transaction_id"`
	GatewayCode          string    `json:"gateway_code"`
	Amount               float64   `json:"amount"`
	Reference            string    `json:"reference"`
	CorrelationID        string    `json:"correlation_id"`
	Timestamp            time.Time `json:"timestamp"`
}

// Bundle es el router centralizado de pagos
type Bundle struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

// New crea e inicia el router de pagos
func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
) *Bundle {
	logger = logger.WithModule("pay.router")

	b := &Bundle{
		rabbit: rabbit,
		log:    logger,
	}

	if rabbit == nil {
		logger.Warn(context.Background()).Msg("RabbitMQ no disponible, pay router deshabilitado")
		return b
	}

	go func() {
		ctx := context.Background()
		logger.Info(ctx).Msg("Starting pay router in background...")
		if err := b.startRouter(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("Pay router failed")
		}
	}()

	logger.Info(context.Background()).Msg("Pay router initialized")
	return b
}

func (b *Bundle) startRouter(ctx context.Context) error {
	if b.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	if err := b.rabbit.DeclareQueue(QueuePayRequests, true); err != nil {
		b.log.Error(ctx).Err(err).Str("queue", QueuePayRequests).Msg("Failed to declare pay.requests queue")
		return err
	}

	// Declarar colas de gateways
	for _, q := range []string{QueueNequiRequests} {
		if err := b.rabbit.DeclareQueue(q, true); err != nil {
			b.log.Warn(ctx).Err(err).Str("queue", q).Msg("Failed to declare gateway queue")
		}
	}

	b.log.Info(ctx).Str("queue", QueuePayRequests).Msg("Pay router listening")
	return b.rabbit.Consume(ctx, QueuePayRequests, b.handlePaymentRequest)
}

func (b *Bundle) handlePaymentRequest(message []byte) error {
	ctx := context.Background()

	var header payRequestHeader
	if err := json.Unmarshal(message, &header); err != nil {
		b.log.Error(ctx).Err(err).Str("body", string(message)).Msg("Failed to unmarshal pay request header")
		return err
	}

	b.log.Info(ctx).
		Uint("transaction_id", header.PaymentTransactionID).
		Str("gateway", header.GatewayCode).
		Msg("Routing payment request")

	targetQueue := b.getGatewayQueue(header.GatewayCode)
	if targetQueue == "" {
		b.log.Error(ctx).Str("gateway", header.GatewayCode).Msg("Unknown payment gateway")
		return nil // No re-encolar si el gateway es desconocido
	}

	if err := b.rabbit.Publish(ctx, targetQueue, message); err != nil {
		b.log.Error(ctx).Err(err).Str("target_queue", targetQueue).Msg("Failed to forward payment request")
		return err
	}

	b.log.Info(ctx).
		Uint("transaction_id", header.PaymentTransactionID).
		Str("gateway", header.GatewayCode).
		Str("target_queue", targetQueue).
		Msg("Payment request forwarded")

	return nil
}

func (b *Bundle) getGatewayQueue(gateway string) string {
	switch gateway {
	case "nequi":
		return QueueNequiRequests
	default:
		return ""
	}
}
