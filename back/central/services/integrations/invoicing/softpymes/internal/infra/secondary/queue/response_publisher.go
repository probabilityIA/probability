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
	QueueInvoiceResponses = "invoicing.responses"
)

// InvoiceResponseMessage es el mensaje que se publica de vuelta a Invoicing Module
type InvoiceResponseMessage struct {
	InvoiceID      uint                   `json:"invoice_id"`
	Provider       string                 `json:"provider"` // "softpymes"
	Status         string                 `json:"status"`   // "success", "error"
	InvoiceNumber  string                 `json:"invoice_number,omitempty"`
	ExternalID     string                 `json:"external_id,omitempty"`
	InvoiceURL     string                 `json:"invoice_url,omitempty"`
	PDFURL         string                 `json:"pdf_url,omitempty"`
	XMLURL         string                 `json:"xml_url,omitempty"`
	CUFE           string                 `json:"cufe,omitempty"`
	IssuedAt       *time.Time             `json:"issued_at,omitempty"`
	DocumentJSON   map[string]interface{} `json:"document_json,omitempty"`
	Error          string                 `json:"error,omitempty"`
	ErrorCode      string                 `json:"error_code,omitempty"`
	ErrorDetails   map[string]interface{} `json:"error_details,omitempty"`
	CorrelationID  string                 `json:"correlation_id"`
	Timestamp      time.Time              `json:"timestamp"`
	ProcessingTime int64                  `json:"processing_time_ms"`

	// Audit data del request/response HTTP al proveedor
	AuditRequestURL     string                 `json:"audit_request_url,omitempty"`
	AuditRequestPayload map[string]interface{} `json:"audit_request_payload,omitempty"`
	AuditResponseStatus int                    `json:"audit_response_status,omitempty"`
	AuditResponseBody   string                 `json:"audit_response_body,omitempty"`
}

// ResponsePublisher publica responses de facturación
type ResponsePublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

// NewResponsePublisher crea un nuevo publisher de responses
func NewResponsePublisher(queue rabbitmq.IQueue, logger log.ILogger) *ResponsePublisher {
	return &ResponsePublisher{
		queue: queue,
		log:   logger.WithModule("softpymes.response_publisher"),
	}
}

// PublishResponse publica una respuesta de facturación
func (p *ResponsePublisher) PublishResponse(ctx context.Context, response *InvoiceResponseMessage) error {
	// Asegurar timestamp
	if response.Timestamp.IsZero() {
		response.Timestamp = time.Now()
	}

	// Serializar mensaje
	data, err := json.Marshal(response)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("Failed to marshal response")
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Publicar en RabbitMQ
	if p.queue == nil {
		p.log.Warn(ctx).
			Uint("invoice_id", response.InvoiceID).
			Msg("RabbitMQ client is nil, cannot publish response")
		return nil // No retornamos error para no romper el flujo, pero logueamos
	}

	if err := p.queue.Publish(ctx, QueueInvoiceResponses, data); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", QueueInvoiceResponses).
			Uint("invoice_id", response.InvoiceID).
			Str("status", response.Status).
			Msg("Failed to publish response")
		return fmt.Errorf("failed to publish response: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", QueueInvoiceResponses).
		Uint("invoice_id", response.InvoiceID).
		Str("status", response.Status).
		Str("correlation_id", response.CorrelationID).
		Int64("processing_time_ms", response.ProcessingTime).
		Msg("📤 Response published successfully")

	return nil
}
