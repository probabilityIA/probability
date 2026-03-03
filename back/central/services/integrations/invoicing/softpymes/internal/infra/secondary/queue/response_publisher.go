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
	QueueInvoiceResponses = rabbitmq.QueueInvoicingResponses
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

// CompareDocumentDetail ítem de un documento de comparación
type CompareDocumentDetail struct {
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
	Quantity string `json:"quantity"`
	Value    string `json:"value"`
	IVA      string `json:"iva"`
}

// CompareDocument documento del proveedor para la comparación
type CompareDocument struct {
	DocumentNumber string                  `json:"document_number"`
	DocumentDate   string                  `json:"document_date"`
	Total          string                  `json:"total"`
	CustomerNit    string                  `json:"customer_nit"`
	CustomerName   string                  `json:"customer_name"`
	Comment        string                  `json:"comment"`
	Prefix         string                  `json:"prefix"`
	Details        []CompareDocumentDetail `json:"details,omitempty"`
}

// CompareResponseMessage mensaje de respuesta de comparación publicado a invoicing.responses
type CompareResponseMessage struct {
	Operation         string            `json:"operation"` // "compare"
	CorrelationID     string            `json:"correlation_id"`
	BusinessID        uint              `json:"business_id"`
	DateFrom          string            `json:"date_from"`
	DateTo            string            `json:"date_to"`
	ProviderDocuments []CompareDocument `json:"provider_documents"`
	Error             string            `json:"error,omitempty"`
	Timestamp         time.Time         `json:"timestamp"`
}

// ResponsePublisher publica responses de facturación
type ResponsePublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

// New crea un nuevo publisher de responses
func New(queue rabbitmq.IQueue, logger log.ILogger) *ResponsePublisher {
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

// PublishCompareResponse publica el resultado de comparación de facturas
func (p *ResponsePublisher) PublishCompareResponse(ctx context.Context, response *CompareResponseMessage) error {
	if response.Timestamp.IsZero() {
		response.Timestamp = time.Now()
	}

	data, err := json.Marshal(response)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("Failed to marshal compare response")
		return fmt.Errorf("failed to marshal compare response: %w", err)
	}

	if p.queue == nil {
		p.log.Warn(ctx).
			Str("correlation_id", response.CorrelationID).
			Msg("RabbitMQ client is nil, cannot publish compare response")
		return nil
	}

	if err := p.queue.Publish(ctx, QueueInvoiceResponses, data); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", QueueInvoiceResponses).
			Str("correlation_id", response.CorrelationID).
			Msg("Failed to publish compare response")
		return fmt.Errorf("failed to publish compare response: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", QueueInvoiceResponses).
		Str("correlation_id", response.CorrelationID).
		Int("documents", len(response.ProviderDocuments)).
		Msg("📤 Compare response published successfully")

	return nil
}
