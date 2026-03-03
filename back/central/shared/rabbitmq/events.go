package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Aliases de retrocompatibilidad — las constantes canónicas están en queues.go.
const (
	EventsExchangeName = ExchangeEvents
	EventsQueueName    = QueueEventsUnified
)

// EventEnvelope es el formato JSON que el consumer de eventos espera.
// Cualquier módulo que publique a ExchangeEvents debe usar este formato.
type EventEnvelope struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Category      string                 `json:"category"`
	BusinessID    uint                   `json:"business_id"`
	IntegrationID uint                   `json:"integration_id,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// PublishEvent serializa un EventEnvelope y lo publica al exchange de eventos.
// El routing key es envelope.Type (e.g. "order.created", "integration.sync.order.created").
func PublishEvent(ctx context.Context, queue IQueue, envelope EventEnvelope) error {
	if queue == nil {
		return nil
	}

	if envelope.ID == "" {
		envelope.ID = uuid.New().String()
	}
	if envelope.Timestamp.IsZero() {
		envelope.Timestamp = time.Now()
	}

	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("error serializando evento: %w", err)
	}

	return queue.PublishToExchange(ctx, EventsExchangeName, envelope.Type, jsonBytes)
}
