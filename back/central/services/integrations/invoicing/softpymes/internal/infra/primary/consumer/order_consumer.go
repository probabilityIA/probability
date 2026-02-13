package consumer

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// OrderConsumer consume eventos de órdenes desde RabbitMQ
// y procesa automáticamente la facturación cuando corresponde
type OrderConsumer struct {
	rabbit  rabbitmq.IQueue
	useCase ports.IInvoiceUseCase
	log     log.ILogger
}

// NewOrderConsumer crea una nueva instancia del consumer de órdenes
func NewOrderConsumer(
	rabbit rabbitmq.IQueue,
	useCase ports.IInvoiceUseCase,
	logger log.ILogger,
) *OrderConsumer {
	return &OrderConsumer{
		rabbit:  rabbit,
		useCase: useCase,
		log:     logger.WithModule("softpymes.consumer"),
	}
}

const (
	QueueOrderEvents = "orders.events.invoicing"
)

// Start inicia el consumer y comienza a procesar mensajes
func (c *OrderConsumer) Start(ctx context.Context) error {
	c.log.Info(ctx).Msg("🚀 Starting Softpymes order consumer")

	// Definir la cola de donde se consumirán los mensajes
	queueName := QueueOrderEvents

	// Iniciar consumo de mensajes con handler
	if err := c.rabbit.Consume(ctx, queueName, c.handleOrderEvent); err != nil {
		c.log.Error(ctx).Err(err).Msg("❌ Failed to start consuming from queue")
		return err
	}

	c.log.Info(ctx).
		Str("queue", queueName).
		Msg("✅ Consumer started successfully")

	return nil
}

// handleOrderEvent procesa un mensaje individual de RabbitMQ
func (c *OrderConsumer) handleOrderEvent(message []byte) error {
	ctx := context.Background()

	// Parsear el mensaje
	var event ports.OrderEventMessage
	if err := json.Unmarshal(message, &event); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("body", string(message)).
			Msg("❌ Failed to unmarshal message")
		return err
	}

	c.log.Info(ctx).
		Str("event_type", event.EventType).
		Str("order_id", event.OrderID).
		Msg("📩 Received order event")

	// Filtrar solo eventos order.created
	if event.EventType != "order.created" {
		c.log.Debug(ctx).
			Str("event_type", event.EventType).
			Msg("⏩ Skipping non-creation event")
		return nil // No es error, simplemente ignorar
	}

	// Validar que el snapshot de orden esté presente
	if event.Order == nil {
		c.log.Warn(ctx).
			Str("event_id", event.EventID).
			Msg("⚠️ Order snapshot is nil - skipping")
		return nil // No es error de procesamiento
	}

	// Procesar la orden para facturación
	if err := c.useCase.ProcessOrderForInvoicing(ctx, &event); err != nil {
		c.log.Error(ctx).
			Err(err).
			Str("order_id", event.OrderID).
			Msg("❌ Failed to process order for invoicing")
		return err // Retornar error para requeue
	}

	c.log.Info(ctx).
		Str("order_id", event.OrderID).
		Msg("✅ Order processed successfully")

	return nil
}
