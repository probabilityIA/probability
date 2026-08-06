package consumerorder

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
		if whaErrors.IsNonRetryable(err) {
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
		if event.BusinessID != nil && *event.BusinessID > 0 {
			trackingURL += "&b=" + strconv.FormatUint(uint64(*event.BusinessID), 10)
		}
	}
	amountToCollect := event.CodTotal
	if amountToCollect <= 0 {
		amountToCollect = event.TotalAmount
	}
	if amountToCollect > 0 && event.CodCarrierFee > 0 {
		amountToCollect += event.CodCarrierFee
	}

	switch templateName {
	case "pedido_en_reparto_cod":
		return map[string]string{
			"1": orDefault(event.CustomerName, "Cliente"),
			"2": orDefault(event.BusinessName, "Probability"),
			"3": orDefault(event.OrderNumber, "N/A"),
			"4": orDefault(event.TrackingNumber, "N/A"),
			"5": orDefault(event.Carrier, "Transportadora"),
			"6": formatTotalAmount(amountToCollect, event.Currency),
			"7": trackingURL,
		}
	case "pedido_en_reparto":
		return map[string]string{
			"1": orDefault(event.CustomerName, "Cliente"),
			"2": orDefault(event.BusinessName, "Probability"),
			"3": orDefault(event.OrderNumber, "N/A"),
			"4": orDefault(event.TrackingNumber, "N/A"),
			"5": orDefault(event.Carrier, "Transportadora"),
			"6": trackingURL,
		}
	case "pedido_entregado_cod":
		return map[string]string{
			"1":  orDefault(event.CustomerName, "Cliente"),
			"2":  orDefault(event.BusinessName, "Probability"),
			"3":  orDefault(event.OrderNumber, "N/A"),
			"4":  orDefault(event.ShippingStreet, orDefault(event.ShippingAddress, "No especificada")),
			"5":  orDefault(event.ShippingCity, ""),
			"6":  orDefault(event.ShippingState, ""),
			"7":  orDefault(event.ItemsSummary, "Ver detalle en plataforma"),
			"8":  orDefault(event.PaymentMethodName, "contra entrega"),
			"9":  orDefault(event.TrackingNumber, "N/A"),
			"10": orDefault(event.Carrier, "Transportadora"),
			"11": formatTotalAmount(amountToCollect, event.Currency),
			"12": trackingURL,
		}
	case "pedido_entregado":
		return map[string]string{
			"1":  orDefault(event.CustomerName, "Cliente"),
			"2":  orDefault(event.BusinessName, "Probability"),
			"3":  orDefault(event.OrderNumber, "N/A"),
			"4":  orDefault(event.ShippingStreet, orDefault(event.ShippingAddress, "No especificada")),
			"5":  orDefault(event.ShippingCity, ""),
			"6":  orDefault(event.ShippingState, ""),
			"7":  orDefault(event.ItemsSummary, "Ver detalle en plataforma"),
			"8":  orDefault(event.TrackingNumber, "N/A"),
			"9":  orDefault(event.Carrier, "Transportadora"),
			"10": trackingURL,
		}
	case "confirmacion_pedido_contraentrega_sin_valor":
		return map[string]string{
			"1": orDefault(event.CustomerName, "Cliente"),
			"2": orDefault(event.BusinessName, "Probability"),
			"3": orDefault(event.OrderNumber, "N/A"),
			"4": orDefault(event.ShippingStreet, orDefault(event.ShippingAddress, "No especificada")),
			"5": orDefault(event.ShippingCity, ""),
			"6": orDefault(event.ShippingState, ""),
			"7": orDefault(event.ItemsSummary, "Ver detalle en plataforma"),
			"8": orDefault(event.PaymentMethodName, "contra entrega"),
		}
	case "confirmacion_pedido":
		return map[string]string{
			"1": orDefault(event.CustomerName, "Cliente"),
			"2": orDefault(event.BusinessName, "Probability"),
			"3": orDefault(event.OrderNumber, "N/A"),
			"4": orDefault(event.ShippingStreet, orDefault(event.ShippingAddress, "No especificada")),
			"5": orDefault(event.ShippingCity, ""),
			"6": orDefault(event.ShippingState, ""),
			"7": orDefault(event.ItemsSummary, "Ver detalle en plataforma"),
		}
	default:
		return map[string]string{
			"1": orDefault(event.CustomerName, "Cliente"),
			"2": orDefault(event.BusinessName, "Probability"),
			"3": orDefault(event.OrderNumber, "N/A"),
			"4": orDefault(event.ShippingStreet, orDefault(event.ShippingAddress, "No especificada")),
			"5": orDefault(event.ShippingCity, ""),
			"6": orDefault(event.ShippingState, ""),
			"7": orDefault(event.ItemsSummary, "Ver detalle en plataforma"),
			"8": orDefault(event.PaymentMethodName, "contra entrega"),
			"9": formatTotalAmount(amountToCollect, event.Currency),
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
