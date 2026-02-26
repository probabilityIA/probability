package queue

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

// RequestPublisher publica solicitudes de pago a la cola pay.requests
type RequestPublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

// New crea un nuevo publisher de solicitudes de pago
func New(queue rabbitmq.IQueue, logger log.ILogger) ports.IRequestPublisher {
	return &RequestPublisher{
		queue: queue,
		log:   logger.WithModule("pay.request_publisher"),
	}
}

// PublishPaymentRequest publica una solicitud de pago
func (p *RequestPublisher) PublishPaymentRequest(ctx context.Context, msg *dtos.PaymentRequestMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal payment request: %w", err)
	}

	if p.queue == nil {
		p.log.Warn(ctx).Uint("transaction_id", msg.PaymentTransactionID).Msg("RabbitMQ not available")
		return fmt.Errorf("rabbitmq not available")
	}

	if err := p.queue.Publish(ctx, constants.QueuePayRequests, data); err != nil {
		p.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to publish payment request")
		return fmt.Errorf("failed to publish payment request: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", constants.QueuePayRequests).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("gateway", msg.GatewayCode).
		Msg("Payment request published")

	return nil
}
