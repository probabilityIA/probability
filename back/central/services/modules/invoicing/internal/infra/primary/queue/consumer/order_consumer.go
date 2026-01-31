package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
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
	QueueOrderEvents = "orders.events"
)

// Start inicia el consumo de eventos de órdenes
func (c *OrderConsumer) Start(ctx context.Context) error {
	c.log.Info(ctx).Str("queue", QueueOrderEvents).Msg("Starting order consumer")

	// Declarar la cola si no existe
	if err := c.queue.DeclareQueue(QueueOrderEvents, true); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to declare queue")
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Iniciar consumo
	if err := c.queue.Consume(ctx, QueueOrderEvents, c.handleOrderEvent); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to start consuming")
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.log.Info(ctx).Msg("Order consumer started successfully")
	return nil
}

// handleOrderEvent procesa un evento de orden
func (c *OrderConsumer) handleOrderEvent(message []byte) error {
	ctx := context.Background()

	c.log.Debug(ctx).Int("size", len(message)).Msg("Received order event")

	// Deserializar evento
	var event request.OrderEvent
	if err := json.Unmarshal(message, &event); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal order event")
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	c.log.Info(ctx).
		Str("event_type", event.EventType).
		Str("order_id", event.OrderID).
		Uint("business_id", event.BusinessID).
		Msg("Processing order event")

	// Solo procesar eventos relevantes
	if !c.shouldProcessEvent(event.EventType) {
		c.log.Debug(ctx).
			Str("event_type", event.EventType).
			Msg("Skipping event - not relevant for invoicing")
		return nil
	}

	// Crear factura automáticamente
	if err := c.createInvoiceForOrder(ctx, &event); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("order_id", event.OrderID).
			Msg("Failed to create invoice for order")
		// No retornar error para no hacer requeue del mensaje
		// El sistema de reintentos de invoicing manejará los fallos
		return nil
	}

	c.log.Info(ctx).
		Str("order_id", event.OrderID).
		Msg("Invoice created successfully from order event")

	return nil
}

// shouldProcessEvent determina si un evento debe procesarse
func (c *OrderConsumer) shouldProcessEvent(eventType string) bool {
	// Eventos relevantes para facturación
	relevantEvents := map[string]bool{
		"order.created": true,
		"order.paid":    true,
		"order.updated": false, // Por ahora no procesamos actualizaciones
	}

	return relevantEvents[eventType]
}

// createInvoiceForOrder crea una factura para una orden
func (c *OrderConsumer) createInvoiceForOrder(ctx context.Context, event *request.OrderEvent) error {
	// Preparar DTO para crear factura
	dto := &dtos.CreateInvoiceDTO{
		OrderID:  event.OrderID,
		IsManual: false, // Es automático desde el consumer
	}

	// Intentar crear la factura
	// El caso de uso manejará:
	// - Validación de configuración (enabled, auto_invoice)
	// - Filtros (min_amount, payment_status, etc.)
	// - Verificación de factura duplicada
	// - Creación de factura
	_, err := c.useCase.CreateInvoice(ctx, dto)
	if err != nil {
		// Log del error pero no fallar el mensaje
		// Esto permite que el sistema continue procesando otros mensajes
		c.log.Warn(ctx).
			Err(err).
			Str("order_id", event.OrderID).
			Msg("Could not create invoice for order - might not meet criteria")
		return err
	}

	return nil
}
