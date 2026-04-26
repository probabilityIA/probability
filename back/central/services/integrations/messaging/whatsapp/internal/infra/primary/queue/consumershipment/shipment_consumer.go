package consumershipment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	whaErrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumershipment/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// Start inicia el consumidor de guías de envío
func (c *consumer) Start(ctx context.Context) error {
	queueName := rabbitmq.QueueShipmentsWhatsAppGuideNotification
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Error declaring shipment guide notification queue")
		return err
	}

	go func() {
		if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
			c.log.Error().Err(err).Msg("Error consuming shipment guide notification queue")
		}
	}()

	return nil
}

// handleMessage procesa cada mensaje de notificación de guía de envío
func (c *consumer) handleMessage(messageBody []byte) error {
	var event request.ShipmentGuideEvent
	if err := json.Unmarshal(messageBody, &event); err != nil {
		c.log.Warn().
			Err(err).
			Msg("Malformed shipment guide message - discarding (ACK)")
		return nil
	}

	c.log.Info().
		Str("tracking_number", event.TrackingNumber).
		Str("order_number", event.OrderNumber).
		Str("customer_phone", event.CustomerPhone).
		Msg("Processing shipment guide notification")

	// Validate phone
	if event.CustomerPhone == "" {
		c.log.Warn().
			Str("order_number", event.OrderNumber).
			Str("tracking_number", event.TrackingNumber).
			Msg("Shipment has no customer phone - skipping notification")
		return nil
	}

	codAmount := event.CodTotal
	if codAmount == 0 {
		codAmount = event.TotalAmount
	}
	trackingURL := event.TrackingURL
	if trackingURL == "" && event.TrackingNumber != "" {
		trackingURL = "https://www.probabilityia.com.co/rastreo?tracking=" + event.TrackingNumber
	}

	variables := map[string]string{
		"1": orDefault(event.CustomerName, "Cliente"),
		"2": orDefault(event.BusinessName, "Probability"),
		"3": orDefault(event.OrderNumber, "N/A"),
		"4": orDefault(event.TrackingNumber, "N/A"),
		"5": orDefault(event.Carrier, "Transportadora"),
		"6": formatTotalAmount(codAmount),
		"7": orDefault(trackingURL, "https://www.probabilityia.com.co/rastreo"),
	}

	businessID := uint(0)
	if event.BusinessID != nil {
		businessID = *event.BusinessID
	}

	templateName := "guia_envio_generada"

	messageID, err := c.useCase.SendTemplate(
		context.Background(),
		templateName,
		event.CustomerPhone,
		variables,
		event.OrderNumber,
		businessID,
	)

	if err != nil {
		if isNonRetryableError(err) {
			c.log.Warn().
				Err(err).
				Str("tracking_number", event.TrackingNumber).
				Str("order_number", event.OrderNumber).
				Str("customer_phone", event.CustomerPhone).
				Msg("WhatsApp shipment notification skipped - non-retryable error (ACK)")
			return nil
		}
		c.log.Error().
			Err(err).
			Str("tracking_number", event.TrackingNumber).
			Str("order_number", event.OrderNumber).
			Str("customer_phone", event.CustomerPhone).
			Msg("Error sending shipment guide template - will be retried")
		return err
	}

	c.log.Info().
		Str("tracking_number", event.TrackingNumber).
		Str("order_number", event.OrderNumber).
		Str("template_name", templateName).
		Str("message_id", messageID).
		Msg("Shipment guide notification sent successfully")

	return nil
}

func orDefault(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func formatTotalAmount(amount float64) string {
	intVal := int64(amount)
	s := fmt.Sprintf("%d", intVal)
	formatted := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			formatted += "."
		}
		formatted += string(c)
	}
	return "$" + formatted
}

func isNonRetryableError(err error) bool {
	var templateNotFound *whaErrors.ErrTemplateNotFound
	if errors.As(err, &templateNotFound) {
		return true
	}
	var missingVar *whaErrors.ErrMissingVariable
	if errors.As(err, &missingVar) {
		return true
	}

	errMsg := err.Error()
	nonRetryablePhrases := []string{
		"key not found",
		"no se encontró integración",
		"credenciales no encontradas",
		"número de teléfono inválido",
		"phone_number_id no encontrado",
		"access_token no encontrado",
		"Required parameter is missing",
		"does not exist in",
		"131008",
		"132001",
		"131009",
	}
	for _, phrase := range nonRetryablePhrases {
		if strings.Contains(errMsg, phrase) {
			return true
		}
	}

	return false
}
