package consumerorder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	whaErrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerorder/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Start inicia el consumidor de órdenes
func (c *consumer) Start(ctx context.Context) error {
	// Declarar cola durable
	queueName := rabbitmq.QueueOrdersConfirmationRequested
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

	return nil
}

// handleMessage procesa cada mensaje de confirmación de orden
func (c *consumer) handleMessage(messageBody []byte) error {
	var event request.OrderConfirmationEvent
	if err := json.Unmarshal(messageBody, &event); err != nil {
		// Mensaje malformado: no tiene sentido reencolar, ACK y descartar
		c.log.Warn().
			Err(err).
			Msg("Malformed order confirmation message - discarding (ACK)")
		return nil
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
		"1": orDefault(event.CustomerName, "Cliente"),
		"2": getBusinessName(event.BusinessID),
		"3": orDefault(event.OrderNumber, "N/A"),
		"4": orDefault(event.ShippingAddress, "No especificada"),
		"5": orDefault(event.ItemsSummary, "Ver detalle en plataforma"),
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
		if isNonRetryableError(err) {
			// Error de configuración/datos: no tiene sentido reencolar.
			// ACK el mensaje para que RabbitMQ no lo reprocese indefinidamente.
			c.log.Warn().
				Err(err).
				Str("order_id", event.OrderID).
				Str("order_number", event.OrderNumber).
				Str("customer_phone", event.CustomerPhone).
				Msg("WhatsApp confirmation skipped - non-retryable error (ACK)")
			return nil
		}
		c.log.Error().
			Err(err).
			Str("order_id", event.OrderID).
			Str("order_number", event.OrderNumber).
			Str("customer_phone", event.CustomerPhone).
			Msg("Error sending confirmation template - will be retried")
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

// orDefault retorna el valor si no está vacío, o el default
func orDefault(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
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

// isNonRetryableError determina si un error no debe provocar reencolar el mensaje.
// Los errores de configuración o datos inválidos nunca se resolverán reintentando:
// reencolarlos causa bucles infinitos en RabbitMQ.
func isNonRetryableError(err error) bool {
	// Errores tipados del dominio: siempre no-retriables
	var templateNotFound *whaErrors.ErrTemplateNotFound
	if errors.As(err, &templateNotFound) {
		return true
	}
	var missingVar *whaErrors.ErrMissingVariable
	if errors.As(err, &missingVar) {
		return true
	}

	// Errores de cache (Redis key not found): la integración no está configurada
	errMsg := err.Error()
	nonRetryablePhrases := []string{
		"key not found",               // Redis cache miss
		"no se encontró integración",  // business sin integración WA
		"credenciales no encontradas", // integration sin credenciales en cache
		"número de teléfono inválido", // teléfono inválido
		"phone_number_id no encontrado", // credenciales incompletas
		"access_token no encontrado",    // credenciales incompletas
		"Required parameter is missing", // Meta: variable vacía en template
		"does not exist in",             // Meta: template no aprobado
		"131008",                        // Meta: parámetro requerido faltante
		"132001",                        // Meta: template no existe
		"131009",                        // Meta: parámetro inválido
	}
	for _, phrase := range nonRetryablePhrases {
		if strings.Contains(errMsg, phrase) {
			return true
		}
	}

	return false
}
