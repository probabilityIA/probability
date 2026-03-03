package dtos

import "time"

// EventEnvelope es el formato de serialización para mensajes RabbitMQ.
// Contiene tags JSON porque vive en la capa de DTOs (no entities).
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
