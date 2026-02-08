package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/events/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/events/internal/infra/secondary/redis"
	"github.com/secamc93/probability/back/central/shared/log"
)

// InvoiceEventConsumer consume eventos de facturación desde Redis y los publica al EventManager
type InvoiceEventConsumer struct {
	subscriber   *redis.InvoiceEventSubscriber
	eventManager domain.IEventPublisher
	logger       log.ILogger
}

// IInvoiceEventConsumer define la interfaz del consumidor de eventos de facturación
type IInvoiceEventConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}

// NewInvoiceEventConsumer crea un nuevo consumidor de eventos de facturación
func NewInvoiceEventConsumer(
	subscriber *redis.InvoiceEventSubscriber,
	eventManager domain.IEventPublisher,
	logger log.ILogger,
) IInvoiceEventConsumer {
	return &InvoiceEventConsumer{
		subscriber:   subscriber,
		eventManager: eventManager,
		logger:       logger,
	}
}

// Start inicia el consumidor
func (c *InvoiceEventConsumer) Start(ctx context.Context) error {
	// Iniciar el suscriptor Redis
	if err := c.subscriber.Start(ctx); err != nil {
		return err
	}

	// Iniciar worker para procesar eventos
	go c.processEvents(ctx)

	return nil
}

// processEvents procesa los eventos recibidos de Redis
func (c *InvoiceEventConsumer) processEvents(ctx context.Context) {
	eventChan := c.subscriber.GetEventChannel()

	for {
		select {
		case event := <-eventChan:
			if event == nil {
				continue
			}

			c.logger.Info(ctx).
				Str("event_id", event.ID).
				Str("event_type", string(event.Type)).
				Uint("business_id", event.BusinessID).
				Msg("Evento de facturación aprobado para notificación, publicando...")

			c.publishInvoiceEvent(ctx, event)

		case <-ctx.Done():
			c.logger.Info(ctx).Msg("Context cancelado, deteniendo procesador de eventos de facturación")
			return
		}
	}
}

// publishInvoiceEvent convierte un InvoiceEvent a Event genérico y lo publica al EventManager
func (c *InvoiceEventConsumer) publishInvoiceEvent(ctx context.Context, invoiceEvent *domain.InvoiceEvent) {
	businessIDStr := fmt.Sprintf("%d", invoiceEvent.BusinessID)

	// Construir metadata
	metadata := map[string]interface{}{
		"business_id": invoiceEvent.BusinessID,
		"event_type":  string(invoiceEvent.Type),
	}

	// Agregar campos relevantes del data al metadata para filtros
	if invoiceID, ok := invoiceEvent.Data["invoice_id"]; ok {
		metadata["invoice_id"] = invoiceID
	}
	if orderID, ok := invoiceEvent.Data["order_id"]; ok {
		metadata["order_id"] = orderID
	}
	if jobID, ok := invoiceEvent.Data["job_id"]; ok {
		metadata["job_id"] = jobID
	}

	genericEvent := domain.Event{
		ID:         invoiceEvent.ID,
		Type:       domain.EventType(invoiceEvent.Type),
		BusinessID: businessIDStr,
		Timestamp:  invoiceEvent.Timestamp,
		Data:       invoiceEvent.Data,
		Metadata:   metadata,
	}

	// Publicar el evento al EventManager (SSE broadcast)
	c.eventManager.PublishEvent(genericEvent)

	c.logger.Info(ctx).
		Str("event_id", invoiceEvent.ID).
		Str("event_type", string(invoiceEvent.Type)).
		Str("business_id", businessIDStr).
		Msg("Evento de facturación publicado al EventManager (SSE)")
}

// Stop detiene el consumidor
func (c *InvoiceEventConsumer) Stop() error {
	return c.subscriber.Stop()
}
