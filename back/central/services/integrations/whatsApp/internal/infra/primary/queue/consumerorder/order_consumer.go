package consumerorder

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/queue/consumerorder/request"
)

// Start inicia el consumidor de órdenes
func (c *consumer) Start(ctx context.Context) error {
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
func (c *consumer) handleMessage(messageBody []byte) error {
	var event request.OrderConfirmationEvent
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

	// Determinar plantilla a usar (desde config o fallback)
	templateName := event.TemplateName
	if templateName == "" {
		// Fallback a plantilla por defecto
		templateName = "confirmacion_pedido_contraentrega"
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

	// Enviar plantilla configurada
	messageID, err := c.useCase.SendTemplate(
		context.Background(),
		templateName, // ← Ahora viene del evento (o fallback)
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
		Str("template_name", templateName).
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
