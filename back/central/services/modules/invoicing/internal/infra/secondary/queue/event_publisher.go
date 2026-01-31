package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// EventPublisher implementa IEventPublisher
type EventPublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

// NewEventPublisher crea un nuevo publicador de eventos
func NewEventPublisher(queue rabbitmq.IQueue, logger log.ILogger) ports.IEventPublisher {
	return &EventPublisher{
		queue: queue,
		log:   logger.WithModule("invoicing.event_publisher"),
	}
}

// InvoiceEvent representa un evento de factura
type InvoiceEvent struct {
	EventType     string  `json:"event_type"`      // "invoice.created", "invoice.cancelled", "invoice.failed"
	InvoiceID     uint    `json:"invoice_id"`
	OrderID       string  `json:"order_id"`
	BusinessID    uint    `json:"business_id"`
	InvoiceNumber string  `json:"invoice_number"`
	TotalAmount   float64 `json:"total_amount"`
	Status        string  `json:"status"`
	ErrorMessage  string  `json:"error_message,omitempty"`
	Timestamp     string  `json:"timestamp"`
}

// CreditNoteEvent representa un evento de nota de crédito
type CreditNoteEvent struct {
	EventType        string  `json:"event_type"` // "credit_note.created"
	CreditNoteID     uint    `json:"credit_note_id"`
	InvoiceID        uint    `json:"invoice_id"`
	BusinessID       uint    `json:"business_id"`
	CreditNoteNumber string  `json:"credit_note_number"`
	Amount           float64 `json:"amount"`
	NoteType         string  `json:"note_type"`
	Timestamp        string  `json:"timestamp"`
}

const (
	QueueInvoiceEvents = "invoicing.events"
)

// PublishInvoiceCreated publica un evento de factura creada
func (p *EventPublisher) PublishInvoiceCreated(ctx context.Context, invoice *entities.Invoice) error {
	event := InvoiceEvent{
		EventType:     "invoice.created",
		InvoiceID:     invoice.ID,
		OrderID:       invoice.OrderID,
		BusinessID:    invoice.BusinessID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
		Status:        invoice.Status,
		Timestamp:     invoice.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return p.publishEvent(ctx, event)
}

// PublishInvoiceCancelled publica un evento de factura cancelada
func (p *EventPublisher) PublishInvoiceCancelled(ctx context.Context, invoice *entities.Invoice) error {
	event := InvoiceEvent{
		EventType:     "invoice.cancelled",
		InvoiceID:     invoice.ID,
		OrderID:       invoice.OrderID,
		BusinessID:    invoice.BusinessID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
		Status:        invoice.Status,
		Timestamp:     invoice.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return p.publishEvent(ctx, event)
}

// PublishInvoiceFailed publica un evento de factura fallida
func (p *EventPublisher) PublishInvoiceFailed(ctx context.Context, invoice *entities.Invoice, errorMsg string) error {
	event := InvoiceEvent{
		EventType:     "invoice.failed",
		InvoiceID:     invoice.ID,
		OrderID:       invoice.OrderID,
		BusinessID:    invoice.BusinessID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
		Status:        invoice.Status,
		ErrorMessage:  errorMsg,
		Timestamp:     invoice.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return p.publishEvent(ctx, event)
}

// PublishCreditNoteCreated publica un evento de nota de crédito creada
func (p *EventPublisher) PublishCreditNoteCreated(ctx context.Context, creditNote *entities.CreditNote) error {
	event := CreditNoteEvent{
		EventType:        "credit_note.created",
		CreditNoteID:     creditNote.ID,
		InvoiceID:        creditNote.InvoiceID,
		BusinessID:       creditNote.BusinessID,
		CreditNoteNumber: creditNote.CreditNoteNumber,
		Amount:           creditNote.Amount,
		NoteType:         creditNote.NoteType,
		Timestamp:        creditNote.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return p.publishEvent(ctx, event)
}

// publishEvent publica un evento genérico
func (p *EventPublisher) publishEvent(ctx context.Context, event interface{}) error {
	// Serializar evento
	data, err := json.Marshal(event)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("Failed to marshal event")
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publicar en RabbitMQ
	if err := p.queue.Publish(ctx, QueueInvoiceEvents, data); err != nil {
		p.log.Error(ctx).Err(err).Str("queue", QueueInvoiceEvents).Msg("Failed to publish event")
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", QueueInvoiceEvents).
		Int("size", len(data)).
		Msg("Event published successfully")

	return nil
}
