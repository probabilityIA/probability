package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/primary/queue/consumer/request"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// OrderConsumer consume eventos de órdenes y crea facturas automáticamente
type OrderConsumer struct {
	queue   rabbitmq.IQueue
	useCase ports.IUseCase
	log     log.ILogger
}

// NewOrderConsumer crea un nuevo consumer de órdenes
func NewOrderConsumer(queue rabbitmq.IQueue, useCase ports.IUseCase, logger log.ILogger) *OrderConsumer {
	return &OrderConsumer{
		queue:   queue,
		useCase: useCase,
		log:     logger.WithModule("invoicing.order_consumer"),
	}
}

const (
	QueueOrderEvents = rabbitmq.QueueOrdersToInvoicing
)

// Start inicia el consumo de eventos de órdenes
func (c *OrderConsumer) Start(ctx context.Context) error {
	// Declarar la cola si no existe
	if err := c.queue.DeclareQueue(QueueOrderEvents, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al declarar cola de facturación")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Iniciar consumo
	if err := c.queue.Consume(ctx, QueueOrderEvents, c.handleOrderEvent); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al iniciar consumer de facturación")
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	return nil
}

// handleOrderEvent procesa un evento de orden
func (c *OrderConsumer) handleOrderEvent(message []byte) error {
	ctx := context.Background()

	// Deserializar evento
	var event request.OrderEvent
	if err := json.Unmarshal(message, &event); err != nil {
		c.log.Error(ctx).Err(err).Msg("Error al deserializar evento de orden")
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Solo procesar eventos relevantes
	if !c.shouldProcessEvent(event.EventType) {
		return nil
	}

	invoice, err := c.useCase.DispatchOrderInvoicing(ctx, event.OrderID)
	if err != nil {
		return nil
	}

	c.log.Info(ctx).
		Str("order_id", event.OrderID).
		Str("invoice_number", invoice.InvoiceNumber).
		Msg("Documento de facturacion creado exitosamente")

	return nil
}

// shouldProcessEvent determina si un evento debe procesarse
func (c *OrderConsumer) shouldProcessEvent(eventType string) bool {
	// Eventos relevantes para facturación
	relevantEvents := map[string]bool{
		"order.created":         true,
		"order.paid":            true,
		"order.guide_generated": true,
		"order.updated":         false,
	}

	return relevantEvents[eventType]
}
