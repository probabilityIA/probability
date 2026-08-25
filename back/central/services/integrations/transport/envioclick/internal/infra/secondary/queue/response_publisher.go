package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	QueueTransportResponses = rabbitmq.QueueTransportResponses
)

type TransportResponseMessage struct {
	ShipmentID    *uint                  `json:"shipment_id,omitempty"`
	BusinessID    uint                   `json:"business_id"`
	Provider      string                 `json:"provider"`
	Operation     string                 `json:"operation"`
	Status        string                 `json:"status"`
	CorrelationID string                 `json:"correlation_id"`
	IsTest        bool                   `json:"is_test,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data,omitempty"`
	Error         string                 `json:"error,omitempty"`
	ErrorKind     string                 `json:"error_kind,omitempty"`
}

type ResponsePublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

func NewResponsePublisher(queue rabbitmq.IQueue, logger log.ILogger) *ResponsePublisher {
	return &ResponsePublisher{
		queue: queue,
		log:   logger.WithModule("transport.envioclick.response_publisher"),
	}
}

func (p *ResponsePublisher) PublishResponse(ctx context.Context, response *TransportResponseMessage) error {
	if response.Timestamp.IsZero() {
		response.Timestamp = time.Now()
	}

	data, err := json.Marshal(response)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("Failed to marshal response")
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	if p.queue == nil {
		p.log.Warn(ctx).
			Str("correlation_id", response.CorrelationID).
			Msg("RabbitMQ client is nil, cannot publish response")
		return nil
	}

	if err := p.queue.Publish(ctx, QueueTransportResponses, data); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", QueueTransportResponses).
			Str("status", response.Status).
			Msg("Failed to publish response")
		return fmt.Errorf("failed to publish response: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", QueueTransportResponses).
		Str("provider", response.Provider).
		Str("operation", response.Operation).
		Str("status", response.Status).
		Str("correlation_id", response.CorrelationID).
		Msg("📤 Transport response published successfully")

	return nil
}
