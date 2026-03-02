package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/queue/mappers"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	QueueInvoiceEvents   = rabbitmq.QueueInvoicingEvents
	QueueBulkInvoiceJobs = rabbitmq.QueueInvoicingBulkCreate
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

// PublishInvoiceCreated publica un evento de factura creada
func (p *EventPublisher) PublishInvoiceCreated(ctx context.Context, invoice *entities.Invoice) error {
	event := mappers.InvoiceToCreatedEvent(invoice)
	if err := p.publishEvent(ctx, QueueInvoiceEvents, event); err != nil {
		return err
	}
	// Notificar al dispatcher central (SSE/WhatsApp/Email)
	go func() {
		_ = rabbitmq.PublishEvent(context.Background(), p.queue, rabbitmq.EventEnvelope{
			Type:       "invoice.created",
			Category:   "invoice",
			BusinessID: invoice.BusinessID,
			Data: map[string]interface{}{
				"invoice_id":     invoice.ID,
				"order_id":       invoice.OrderID,
				"invoice_number": invoice.InvoiceNumber,
				"total_amount":   invoice.TotalAmount,
				"status":         invoice.Status,
			},
		})
	}()
	return nil
}

// PublishInvoiceCancelled publica un evento de factura cancelada
func (p *EventPublisher) PublishInvoiceCancelled(ctx context.Context, invoice *entities.Invoice) error {
	event := mappers.InvoiceToCancelledEvent(invoice)
	if err := p.publishEvent(ctx, QueueInvoiceEvents, event); err != nil {
		return err
	}
	go func() {
		_ = rabbitmq.PublishEvent(context.Background(), p.queue, rabbitmq.EventEnvelope{
			Type:       "invoice.cancelled",
			Category:   "invoice",
			BusinessID: invoice.BusinessID,
			Data: map[string]interface{}{
				"invoice_id":     invoice.ID,
				"order_id":       invoice.OrderID,
				"invoice_number": invoice.InvoiceNumber,
				"status":         invoice.Status,
			},
		})
	}()
	return nil
}

// PublishInvoiceFailed publica un evento de factura fallida
func (p *EventPublisher) PublishInvoiceFailed(ctx context.Context, invoice *entities.Invoice, errorMsg string) error {
	event := mappers.InvoiceToFailedEvent(invoice, errorMsg)
	if err := p.publishEvent(ctx, QueueInvoiceEvents, event); err != nil {
		return err
	}
	go func() {
		_ = rabbitmq.PublishEvent(context.Background(), p.queue, rabbitmq.EventEnvelope{
			Type:       "invoice.failed",
			Category:   "invoice",
			BusinessID: invoice.BusinessID,
			Data: map[string]interface{}{
				"invoice_id":     invoice.ID,
				"order_id":       invoice.OrderID,
				"invoice_number": invoice.InvoiceNumber,
				"total_amount":   invoice.TotalAmount,
				"status":         invoice.Status,
				"error_message":  errorMsg,
			},
		})
	}()
	return nil
}

// PublishCreditNoteCreated publica un evento de nota de crédito creada
func (p *EventPublisher) PublishCreditNoteCreated(ctx context.Context, creditNote *entities.CreditNote) error {
	event := mappers.CreditNoteToEvent(creditNote)
	if err := p.publishEvent(ctx, QueueInvoiceEvents, event); err != nil {
		return err
	}
	go func() {
		_ = rabbitmq.PublishEvent(context.Background(), p.queue, rabbitmq.EventEnvelope{
			Type:       "credit_note.created",
			Category:   "invoice",
			BusinessID: creditNote.BusinessID,
			Data: map[string]interface{}{
				"credit_note_id":     creditNote.ID,
				"invoice_id":         creditNote.InvoiceID,
				"credit_note_number": creditNote.CreditNoteNumber,
				"amount":             creditNote.Amount,
				"note_type":          creditNote.NoteType,
			},
		})
	}()
	return nil
}

// PublishBulkInvoiceJob publica un mensaje para procesar una factura en un job masivo
func (p *EventPublisher) PublishBulkInvoiceJob(ctx context.Context, dto *dtos.BulkInvoiceJobMessage) error {
	// Mapear DTO de dominio a mensaje de queue
	message := mappers.BulkJobDTOToMessage(dto)

	// Serializar mensaje
	data, err := json.Marshal(message)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("Failed to marshal bulk invoice job message")
		return fmt.Errorf("failed to marshal bulk invoice job message: %w", err)
	}

	// Publicar en RabbitMQ
	if err := p.queue.Publish(ctx, QueueBulkInvoiceJobs, data); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", QueueBulkInvoiceJobs).
			Str("job_id", dto.JobID).
			Str("order_id", dto.OrderID).
			Msg("Failed to publish bulk invoice job message")
		return fmt.Errorf("failed to publish bulk invoice job message: %w", err)
	}

	p.log.Debug(ctx).
		Str("queue", QueueBulkInvoiceJobs).
		Str("job_id", dto.JobID).
		Str("order_id", dto.OrderID).
		Msg("Bulk invoice job message published")

	return nil
}

// publishEvent publica un evento genérico
func (p *EventPublisher) publishEvent(ctx context.Context, queueName string, event interface{}) error {
	// Serializar evento
	data, err := json.Marshal(event)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("Failed to marshal event")
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publicar en RabbitMQ
	if err := p.queue.Publish(ctx, queueName, data); err != nil {
		p.log.Error(ctx).Err(err).Str("queue", queueName).Msg("Failed to publish event")
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", queueName).
		Int("size", len(data)).
		Msg("Event published successfully")

	return nil
}
