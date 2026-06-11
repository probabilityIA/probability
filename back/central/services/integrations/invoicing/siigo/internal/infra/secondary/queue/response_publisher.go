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

type InvoiceResponseMessage struct {
	InvoiceID      uint                   `json:"invoice_id"`
	Provider       string                 `json:"provider"`
	Status         string                 `json:"status"`
	Operation      string                 `json:"operation,omitempty"`
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

	AuditRequestURL     string                 `json:"audit_request_url,omitempty"`
	AuditRequestPayload map[string]interface{} `json:"audit_request_payload,omitempty"`
	AuditResponseStatus int                    `json:"audit_response_status,omitempty"`
	AuditResponseBody   string                 `json:"audit_response_body,omitempty"`

	CashReceiptRequestURL     string                 `json:"cash_receipt_request_url,omitempty"`
	CashReceiptRequestPayload map[string]interface{} `json:"cash_receipt_request_payload,omitempty"`
	CashReceiptResponseStatus int                    `json:"cash_receipt_response_status,omitempty"`
	CashReceiptResponseBody   string                 `json:"cash_receipt_response_body,omitempty"`
}

type CompareDocumentDetail struct {
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
	Quantity string `json:"quantity"`
	Value    string `json:"value"`
	IVA      string `json:"iva"`
}

type CompareDocument struct {
	DocumentNumber     string                  `json:"document_number"`
	DocumentDate       string                  `json:"document_date"`
	Total              string                  `json:"total"`
	CustomerNit        string                  `json:"customer_nit"`
	CustomerName       string                  `json:"customer_name"`
	Comment            string                  `json:"comment"`
	Prefix             string                  `json:"prefix"`
	Annuled            bool                    `json:"annuled"`
	ElectronicDocument bool                    `json:"electronic_document"`
	Details            []CompareDocumentDetail `json:"details,omitempty"`
}

type CompareResponseMessage struct {
	Operation         string            `json:"operation"`
	Mode              string            `json:"mode,omitempty"`
	CorrelationID     string            `json:"correlation_id"`
	BusinessID        uint              `json:"business_id"`
	DateFrom          string            `json:"date_from"`
	DateTo            string            `json:"date_to"`
	ProviderDocuments []CompareDocument `json:"provider_documents"`
	Error             string            `json:"error,omitempty"`
	Timestamp         time.Time         `json:"timestamp"`
}

type BankAccountItem struct {
	AccountNumber string `json:"account_number"`
	Name          string `json:"name"`
	NameType      string `json:"name_type"`
}

type ListBankAccountsResponseMessage struct {
	Operation     string            `json:"operation"`
	CorrelationID string            `json:"correlation_id"`
	BusinessID    uint              `json:"business_id"`
	Items         []BankAccountItem `json:"items"`
	Error         string            `json:"error,omitempty"`
	Timestamp     time.Time         `json:"timestamp"`
}

type ListItemsItem struct {
	ItemCode      string  `json:"item_code"`
	ItemName      string  `json:"item_name"`
	ItemPrice     float64 `json:"item_price"`
	UnitCost      float64 `json:"unit_cost"`
	Description   string  `json:"description"`
	MinimumStock  string  `json:"minimum_stock"`
	OrderQuantity string  `json:"order_quantity"`
}

type ListItemsResponseMessage struct {
	Operation     string          `json:"operation"`
	CorrelationID string          `json:"correlation_id"`
	BusinessID    uint            `json:"business_id"`
	Items         []ListItemsItem `json:"items"`
	Error         string          `json:"error,omitempty"`
	Timestamp     time.Time       `json:"timestamp"`
}

type ResponsePublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

func New(queue rabbitmq.IQueue, logger log.ILogger) *ResponsePublisher {
	return &ResponsePublisher{
		queue: queue,
		log:   logger.WithModule("siigo.response_publisher"),
	}
}

func (p *ResponsePublisher) publish(ctx context.Context, correlationID string, payload interface{}, kind string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		p.log.Error(ctx).Err(err).Str("kind", kind).Msg("Failed to marshal response")
		return fmt.Errorf("failed to marshal %s response: %w", kind, err)
	}

	if p.queue == nil {
		p.log.Warn(ctx).
			Str("correlation_id", correlationID).
			Str("kind", kind).
			Msg("RabbitMQ client is nil, cannot publish response")
		return nil
	}

	if err := p.queue.Publish(ctx, QueueInvoiceResponses, data); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", QueueInvoiceResponses).
			Str("correlation_id", correlationID).
			Str("kind", kind).
			Msg("Failed to publish response")
		return fmt.Errorf("failed to publish %s response: %w", kind, err)
	}

	p.log.Info(ctx).
		Str("queue", QueueInvoiceResponses).
		Str("correlation_id", correlationID).
		Str("kind", kind).
		Msg("Siigo response published successfully")

	return nil
}

func (p *ResponsePublisher) PublishResponse(ctx context.Context, response *InvoiceResponseMessage) error {
	if response.Timestamp.IsZero() {
		response.Timestamp = time.Now()
	}
	return p.publish(ctx, response.CorrelationID, response, "invoice")
}

func (p *ResponsePublisher) PublishCompareResponse(ctx context.Context, response *CompareResponseMessage) error {
	if response.Timestamp.IsZero() {
		response.Timestamp = time.Now()
	}
	return p.publish(ctx, response.CorrelationID, response, "compare")
}

func (p *ResponsePublisher) PublishListItemsResponse(ctx context.Context, response *ListItemsResponseMessage) error {
	if response.Timestamp.IsZero() {
		response.Timestamp = time.Now()
	}
	return p.publish(ctx, response.CorrelationID, response, "list_items")
}

func (p *ResponsePublisher) PublishListBankAccountsResponse(ctx context.Context, response *ListBankAccountsResponseMessage) error {
	if response.Timestamp.IsZero() {
		response.Timestamp = time.Now()
	}
	return p.publish(ctx, response.CorrelationID, response, "list_bank_accounts")
}
