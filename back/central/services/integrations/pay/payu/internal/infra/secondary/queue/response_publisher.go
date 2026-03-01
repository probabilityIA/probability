package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/pay/payu/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const QueuePayResponses = rabbitmq.QueuePayResponses

// ResponsePublisher publica respuestas de pago PayU
type ResponsePublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

// New crea un nuevo publisher de respuestas
func New(queue rabbitmq.IQueue, logger log.ILogger) ports.IResponsePublisher {
	return &ResponsePublisher{
		queue: queue,
		log:   logger.WithModule("payu.response_publisher"),
	}
}

// PublishPaymentResponse publica una respuesta de pago a pay.responses
func (p *ResponsePublisher) PublishPaymentResponse(ctx context.Context, msg *ports.PaymentResponseMsg) error {
	type msgWithTimestamp struct {
		*ports.PaymentResponseMsg
		Timestamp time.Time `json:"timestamp"`
	}
	enriched := msgWithTimestamp{
		PaymentResponseMsg: msg,
		Timestamp:          time.Now(),
	}

	data, err := json.Marshal(enriched)
	if err != nil {
		return fmt.Errorf("failed to marshal payment response: %w", err)
	}

	if p.queue == nil {
		p.log.Warn(ctx).Uint("transaction_id", msg.PaymentTransactionID).Msg("RabbitMQ not available")
		return nil
	}

	if err := p.queue.Publish(ctx, QueuePayResponses, data); err != nil {
		p.log.Error(ctx).Err(err).Uint("transaction_id", msg.PaymentTransactionID).Msg("Failed to publish payment response")
		return fmt.Errorf("failed to publish payment response: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", QueuePayResponses).
		Uint("transaction_id", msg.PaymentTransactionID).
		Str("status", msg.Status).
		Msg("Payment response published")

	return nil
}
