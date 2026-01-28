package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/app"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// OrderConfirmationEvent representa el evento de confirmación de orden
type OrderConfirmationEvent struct {
	EventType       string  `json:"event_type"`
	OrderID         string  `json:"order_id"`
	OrderNumber     string  `json:"order_number"`
	BusinessID      *uint   `json:"business_id"`
	CustomerName    string  `json:"customer_name"`
	CustomerPhone   string  `json:"customer_phone"`
	CustomerEmail   string  `json:"customer_email"`
	TotalAmount     float64 `json:"total_amount"`
	Currency        string  `json:"currency"`
	ItemsSummary    string  `json:"items_summary"`
	ShippingAddress string  `json:"shipping_address"`
	PaymentMethodID uint    `json:"payment_method_id"`
	IntegrationID   uint    `json:"integration_id"`
	Platform        string  `json:"platform"`
	Timestamp       int64   `json:"timestamp"`
}

// OrderConfirmationConsumer consume eventos de órdenes que requieren confirmación
type OrderConfirmationConsumer struct {
	queue               rabbitmq.IQueue
	sendTemplateUseCase app.ISendTemplateMessageUseCase
	log                 log.ILogger
}

// NewOrderConfirmationConsumer crea un nuevo consumidor de órdenes
func NewOrderConfirmationConsumer(
	queue rabbitmq.IQueue,
	sendTemplateUseCase app.ISendTemplateMessageUseCase,
	logger log.ILogger,
) *OrderConfirmationConsumer {
	return &OrderConfirmationConsumer{
		queue:               queue,
		sendTemplateUseCase: sendTemplateUseCase,
		log:                 logger,
	}
}

// Start inicia el consumidor de órdenes
func (c *OrderConfirmationConsumer) Start(ctx context.Context) error {
	c.log.Info().Msg("Starting order confirmation consumer")

	// Declarar cola durable
	queueName := "orders.confirmation.requested"
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Error declaring queue")
		return err
	}

	// Consumir mensajes
	go func() {
		if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
			c.log.Error().Err(err).Msg("Error consuming order confirmation queue")
		}
	}()

	c.log.Info().
		Str("queue", queueName).
		Msg("Order confirmation consumer started successfully")

	return nil
}

// handleMessage procesa cada mensaje de confirmación de orden
func (c *OrderConfirmationConsumer) handleMessage(messageBody []byte) error {
	var event OrderConfirmationEvent
	if err := json.Unmarshal(messageBody, &event); err != nil {
		c.log.Error().
			Err(err).
			Msg("Error unmarshaling order confirmation event")
		return err
	}

	c.log.Info().
		Str("order_id", event.OrderID).
		Str("order_number", event.OrderNumber).
		Str("customer_phone", event.CustomerPhone).
		Msg("Processing order confirmation request")

	// Validar que tenga teléfono
	if event.CustomerPhone == "" {
		c.log.Warn().
			Str("order_id", event.OrderID).
			Str("order_number", event.OrderNumber).
			Msg("Order has no customer phone - skipping confirmation")
		return nil // No error, solo skip
	}

	// Construir variables para la plantilla
	// Según el plan, la plantilla "confirmacion_pedido_contraentrega" usa:
	// 1: customer_name
	// 2: business_name
	// 3: order_number
	// 4: shipping_address
	// 5: items_summary
	variables := map[string]string{
		"1": event.CustomerName,
		"2": getBusinessName(event.BusinessID), // Helper para obtener nombre del negocio
		"3": event.OrderNumber,
		"4": event.ShippingAddress,
		"5": event.ItemsSummary,
	}

	// Obtener BusinessID (puede ser nulo)
	businessID := uint(0)
	if event.BusinessID != nil {
		businessID = *event.BusinessID
	}

	// Enviar plantilla de confirmación
	messageID, err := c.sendTemplateUseCase.SendTemplate(
		context.Background(),
		"confirmacion_pedido_contraentrega",
		event.CustomerPhone,
		variables,
		event.OrderNumber,
		businessID,
	)

	if err != nil {
		c.log.Error().
			Err(err).
			Str("order_id", event.OrderID).
			Str("order_number", event.OrderNumber).
			Str("customer_phone", event.CustomerPhone).
			Msg("Error sending confirmation template")
		return err
	}

	c.log.Info().
		Str("order_id", event.OrderID).
		Str("order_number", event.OrderNumber).
		Str("message_id", messageID).
		Msg("Confirmation template sent successfully")

	return nil
}

// getBusinessName obtiene el nombre del negocio (implementación simplificada)
// TODO: En producción, esto debería consultar la BD o tener un caché
func getBusinessName(businessID *uint) string {
	if businessID == nil {
		return "Probability"
	}
	// Por ahora retornamos un placeholder
	// En implementación completa, consultar repositorio de Business
	return fmt.Sprintf("Negocio #%d", *businessID)
}
