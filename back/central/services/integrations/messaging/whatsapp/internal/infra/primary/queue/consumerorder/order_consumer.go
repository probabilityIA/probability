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

func (c *consumer) Start(ctx context.Context) error {
	queueName := rabbitmq.QueueOrdersConfirmationRequested
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().
			Err(err).
			Str("queue", queueName).
			Msg("Error declaring queue")
		return err
	}

	go func() {
		if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
			c.log.Error().Err(err).Msg("Error consuming order confirmation queue")
		}
	}()

	return nil
}

func (c *consumer) handleMessage(messageBody []byte) error {
	var event request.OrderConfirmationEvent
	if err := json.Unmarshal(messageBody, &event); err != nil {
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

	if event.CustomerPhone == "" {
		c.log.Warn().
			Str("order_id", event.OrderID).
			Str("order_number", event.OrderNumber).
			Msg("Order has no customer phone - skipping confirmation")
		return nil
	}

	templateName := event.TemplateName
	if templateName == "" {
		templateName = "confirmacion_pedido_contraentrega"
	}

	variables := buildVariables(templateName, event)

	businessID := uint(0)
	if event.BusinessID != nil {
		businessID = *event.BusinessID
	}

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

func buildVariables(templateName string, event request.OrderConfirmationEvent) map[string]string {
	trackingURL := "https://www.probabilityia.com.co/rastreo"
	if event.TrackingNumber != "" {
		trackingURL = "https://www.probabilityia.com.co/rastreo?tracking=" + event.TrackingNumber
	}
	switch templateName {
	case "pedido_en_reparto":
		return map[string]string{
			"1": orDefault(event.CustomerName, "Cliente"),
			"2": orDefault(event.BusinessName, "Probability"),
			"3": orDefault(event.OrderNumber, "N/A"),
			"4": orDefault(event.TrackingNumber, "N/A"),
			"5": orDefault(event.Carrier, "Transportadora"),
			"6": formatTotalAmount(event.TotalAmount, event.Currency),
			"7": trackingURL,
		}
	case "pedido_entregado":
		return map[string]string{
			"1":  orDefault(event.CustomerName, "Cliente"),
			"2":  orDefault(event.BusinessName, "Probability"),
			"3":  orDefault(event.OrderNumber, "N/A"),
			"4":  orDefault(event.ShippingAddress, "No especificada"),
			"5":  orDefault(event.ShippingCity, ""),
			"6":  orDefault(event.ShippingState, ""),
			"7":  orDefault(event.ItemsSummary, "Ver detalle en plataforma"),
			"8":  orDefault(event.PaymentMethodName, "contra entrega"),
			"9":  orDefault(event.TrackingNumber, "N/A"),
			"10": orDefault(event.Carrier, "Transportadora"),
			"11": formatTotalAmount(event.TotalAmount, event.Currency),
			"12": trackingURL,
		}
	default:
		return map[string]string{
			"1": orDefault(event.CustomerName, "Cliente"),
			"2": orDefault(event.BusinessName, "Probability"),
			"3": orDefault(event.OrderNumber, "N/A"),
			"4": orDefault(event.ShippingAddress, "No especificada"),
			"5": orDefault(event.ShippingCity, ""),
			"6": orDefault(event.ShippingState, ""),
			"7": orDefault(event.ItemsSummary, "Ver detalle en plataforma"),
			"8": orDefault(event.PaymentMethodName, "contra entrega"),
			"9": formatTotalAmount(event.TotalAmount, event.Currency),
		}
	}
}

func orDefault(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func formatTotalAmount(amount float64, _ string) string {
	intVal := int64(amount)
	formatted := ""
	s := fmt.Sprintf("%d", intVal)
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
