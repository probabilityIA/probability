package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	QueueTransportRequests = "transport.requests"
)

// TransportRequestPublisher publishes transport requests to RabbitMQ
type TransportRequestPublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

// NewTransportRequestPublisher creates a new publisher
func NewTransportRequestPublisher(queue rabbitmq.IQueue, logger log.ILogger) domain.ITransportRequestPublisher {
	return &TransportRequestPublisher{
		queue: queue,
		log:   logger.WithModule("shipments.transport_request_publisher"),
	}
}

// PublishTransportRequest publishes a transport request to the unified queue
func (p *TransportRequestPublisher) PublishTransportRequest(ctx context.Context, request *domain.TransportRequestMessage) error {
	if request.Timestamp.IsZero() {
		request.Timestamp = time.Now()
	}

	data, err := json.Marshal(request)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("Failed to marshal transport request")
		return fmt.Errorf("failed to marshal transport request: %w", err)
	}

	if p.queue == nil {
		p.log.Warn(ctx).
			Str("correlation_id", request.CorrelationID).
			Msg("RabbitMQ client is nil, cannot publish transport request")
		return fmt.Errorf("rabbitmq client is nil")
	}

	if err := p.queue.Publish(ctx, QueueTransportRequests, data); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", QueueTransportRequests).
			Str("provider", request.Provider).
			Str("operation", request.Operation).
			Str("correlation_id", request.CorrelationID).
			Msg("Failed to publish transport request")
		return fmt.Errorf("failed to publish transport request: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", QueueTransportRequests).
		Str("provider", request.Provider).
		Str("operation", request.Operation).
		Str("correlation_id", request.CorrelationID).
		Msg("📤 Transport request published successfully")

	return nil
}
