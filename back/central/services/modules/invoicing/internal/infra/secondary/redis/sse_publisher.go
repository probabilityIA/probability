package redis

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	redisclient "github.com/secamc93/probability/back/central/shared/redis"
)

// invoiceSSEEvent es la estructura del evento publicado a Redis
type invoiceSSEEvent struct {
	ID         string                 `json:"id"`
	EventType  string                 `json:"event_type"`
	BusinessID uint                   `json:"business_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}

// SSEPublisher publica eventos de facturación a Redis Pub/Sub para SSE
type SSEPublisher struct {
	redisClient redisclient.IRedis
	logger      log.ILogger
	channel     string
}

// NewSSEPublisher crea un nuevo publicador SSE de facturación
func NewSSEPublisher(redisClient redisclient.IRedis, logger log.ILogger, channel string) ports.IInvoiceSSEPublisher {
	return &SSEPublisher{
		redisClient: redisClient,
		logger:      logger,
		channel:     channel,
	}
}

// PublishInvoiceCreated publica evento de factura creada
func (p *SSEPublisher) PublishInvoiceCreated(ctx context.Context, invoice *entities.Invoice) error {
	event := invoiceSSEEvent{
		ID:         generateEventID(),
		EventType:  "invoice.created",
		BusinessID: invoice.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"invoice_id":     invoice.ID,
			"order_id":       invoice.OrderID,
			"invoice_number": invoice.InvoiceNumber,
			"total_amount":   invoice.TotalAmount,
			"currency":       invoice.Currency,
			"status":         invoice.Status,
			"customer_name":  invoice.CustomerName,
			"external_url":   ptrToString(invoice.InvoiceURL),
		},
	}
	return p.publish(ctx, event)
}

// PublishInvoiceFailed publica evento de factura fallida
func (p *SSEPublisher) PublishInvoiceFailed(ctx context.Context, invoice *entities.Invoice, errorMsg string) error {
	event := invoiceSSEEvent{
		ID:         generateEventID(),
		EventType:  "invoice.failed",
		BusinessID: invoice.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"invoice_id":    invoice.ID,
			"order_id":      invoice.OrderID,
			"total_amount":  invoice.TotalAmount,
			"currency":      invoice.Currency,
			"status":        invoice.Status,
			"customer_name": invoice.CustomerName,
			"error_message": errorMsg,
		},
	}
	return p.publish(ctx, event)
}

// PublishInvoiceCancelled publica evento de factura cancelada
func (p *SSEPublisher) PublishInvoiceCancelled(ctx context.Context, invoice *entities.Invoice) error {
	event := invoiceSSEEvent{
		ID:         generateEventID(),
		EventType:  "invoice.cancelled",
		BusinessID: invoice.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"invoice_id":     invoice.ID,
			"order_id":       invoice.OrderID,
			"invoice_number": invoice.InvoiceNumber,
			"total_amount":   invoice.TotalAmount,
			"currency":       invoice.Currency,
			"status":         invoice.Status,
		},
	}
	return p.publish(ctx, event)
}

// PublishCreditNoteCreated publica evento de nota de crédito creada
func (p *SSEPublisher) PublishCreditNoteCreated(ctx context.Context, creditNote *entities.CreditNote) error {
	event := invoiceSSEEvent{
		ID:         generateEventID(),
		EventType:  "credit_note.created",
		BusinessID: creditNote.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"credit_note_id":     creditNote.ID,
			"invoice_id":         creditNote.InvoiceID,
			"credit_note_number": creditNote.CreditNoteNumber,
			"amount":             creditNote.Amount,
			"currency":           creditNote.Currency,
			"status":             creditNote.Status,
			"reason":             creditNote.Reason,
		},
	}
	return p.publish(ctx, event)
}

// PublishBulkJobProgress publica progreso de un job masivo
func (p *SSEPublisher) PublishBulkJobProgress(ctx context.Context, job *entities.BulkInvoiceJob) error {
	event := invoiceSSEEvent{
		ID:         generateEventID(),
		EventType:  "bulk_job.progress",
		BusinessID: job.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"job_id":       job.ID,
			"total_orders": job.TotalOrders,
			"processed":    job.Processed,
			"successful":   job.Successful,
			"failed":       job.Failed,
			"progress":     job.GetProgress(),
			"status":       job.Status,
		},
	}
	return p.publish(ctx, event)
}

// PublishBulkJobCompleted publica que un job masivo finalizó
func (p *SSEPublisher) PublishBulkJobCompleted(ctx context.Context, job *entities.BulkInvoiceJob) error {
	event := invoiceSSEEvent{
		ID:         generateEventID(),
		EventType:  "bulk_job.completed",
		BusinessID: job.BusinessID,
		Timestamp:  time.Now(),
		Data: map[string]interface{}{
			"job_id":       job.ID,
			"total_orders": job.TotalOrders,
			"processed":    job.Processed,
			"successful":   job.Successful,
			"failed":       job.Failed,
			"progress":     100,
			"status":       job.Status,
		},
	}
	return p.publish(ctx, event)
}

// publish serializa y publica el evento a Redis de forma no-bloqueante
func (p *SSEPublisher) publish(ctx context.Context, event invoiceSSEEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Error(ctx).
			Err(err).
			Str("event_type", event.EventType).
			Msg("Error serializing invoice SSE event")
		return err
	}

	// Publicar de forma no-bloqueante
	go func() {
		publishCtx := context.Background()
		if pubErr := p.redisClient.Client(publishCtx).Publish(publishCtx, p.channel, eventJSON).Err(); pubErr != nil {
			p.logger.Error(publishCtx).
				Err(pubErr).
				Str("event_type", event.EventType).
				Str("channel", p.channel).
				Msg("Error publishing invoice SSE event to Redis")
			return
		}

		p.logger.Info(publishCtx).
			Str("event_id", event.ID).
			Str("event_type", event.EventType).
			Uint("business_id", event.BusinessID).
			Str("channel", p.channel).
			Msg("Invoice SSE event published to Redis")
	}()

	return nil
}

// ptrToString convierte un puntero de string a string vacío si es nil
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// generateEventID genera un ID único para el evento
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString genera una cadena aleatoria
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// noopSSEPublisher es una implementación no-op para cuando Redis no está disponible
type noopSSEPublisher struct{}

// NewNoopSSEPublisher crea un publisher que no hace nada (para cuando Redis no está disponible)
func NewNoopSSEPublisher() ports.IInvoiceSSEPublisher {
	return &noopSSEPublisher{}
}

func (n *noopSSEPublisher) PublishInvoiceCreated(_ context.Context, _ *entities.Invoice) error {
	return nil
}
func (n *noopSSEPublisher) PublishInvoiceFailed(_ context.Context, _ *entities.Invoice, _ string) error {
	return nil
}
func (n *noopSSEPublisher) PublishInvoiceCancelled(_ context.Context, _ *entities.Invoice) error {
	return nil
}
func (n *noopSSEPublisher) PublishCreditNoteCreated(_ context.Context, _ *entities.CreditNote) error {
	return nil
}
func (n *noopSSEPublisher) PublishBulkJobProgress(_ context.Context, _ *entities.BulkInvoiceJob) error {
	return nil
}
func (n *noopSSEPublisher) PublishBulkJobCompleted(_ context.Context, _ *entities.BulkInvoiceJob) error {
	return nil
}

// Compile-time interface checks
var _ ports.IInvoiceSSEPublisher = (*SSEPublisher)(nil)
var _ ports.IInvoiceSSEPublisher = (*noopSSEPublisher)(nil)
