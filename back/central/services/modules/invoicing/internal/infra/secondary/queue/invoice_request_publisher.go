package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Colas de requests por proveedor
const (
	QueueSoftpymesRequests = "invoicing.softpymes.requests"
	QueueSiigoRequests     = "invoicing.siigo.requests"
	QueueFactusRequests    = "invoicing.factus.requests"
)

// InvoiceRequestPublisher implementa IInvoiceRequestPublisher
type InvoiceRequestPublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

// NewInvoiceRequestPublisher crea un nuevo publicador de requests de facturación
func NewInvoiceRequestPublisher(queue rabbitmq.IQueue, logger log.ILogger) ports.IInvoiceRequestPublisher {
	return &InvoiceRequestPublisher{
		queue: queue,
		log:   logger.WithModule("invoicing.request_publisher"),
	}
}

// PublishInvoiceRequest publica una solicitud de facturación a la cola del proveedor correspondiente
func (p *InvoiceRequestPublisher) PublishInvoiceRequest(ctx context.Context, request *dtos.InvoiceRequestMessage) error {
	// Determinar cola según proveedor
	queueName := p.getQueueNameForProvider(request.Provider)

	// Serializar request
	data, err := json.Marshal(request)
	if err != nil {
		p.log.Error(ctx).Err(err).Msg("Failed to marshal invoice request")
		return fmt.Errorf("failed to marshal invoice request: %w", err)
	}

	// Publicar en RabbitMQ
	if err := p.queue.Publish(ctx, queueName, data); err != nil {
		p.log.Error(ctx).
			Err(err).
			Str("queue", queueName).
			Str("provider", request.Provider).
			Uint("invoice_id", request.InvoiceID).
			Str("correlation_id", request.CorrelationID).
			Msg("Failed to publish invoice request")
		return fmt.Errorf("failed to publish invoice request: %w", err)
	}

	p.log.Info(ctx).
		Str("queue", queueName).
		Str("provider", request.Provider).
		Str("operation", request.Operation).
		Uint("invoice_id", request.InvoiceID).
		Str("correlation_id", request.CorrelationID).
		Int("size", len(data)).
		Msg("📤 Invoice request published successfully")

	return nil
}

// getQueueNameForProvider retorna el nombre de la cola según el proveedor
func (p *InvoiceRequestPublisher) getQueueNameForProvider(provider string) string {
	switch provider {
	case dtos.ProviderSoftpymes:
		return QueueSoftpymesRequests
	case dtos.ProviderSiigo:
		return QueueSiigoRequests
	case dtos.ProviderFactus:
		return QueueFactusRequests
	default:
		// Default a softpymes para mantener compatibilidad
		return QueueSoftpymesRequests
	}
}
