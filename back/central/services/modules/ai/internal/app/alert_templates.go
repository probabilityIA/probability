package app

import (
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

type statusPhrase struct {
	title    string
	body     string
	severity string
}

var orderStatusPhrases = map[string]statusPhrase{
	"cancelled":        {"Orden cancelada", "fue cancelada", entities.AlertSeverityWarning},
	"failed":           {"Orden fallida", "pas\u00f3 a estado fallido", entities.AlertSeverityCritical},
	"rejected":         {"Orden rechazada", "fue rechazada", entities.AlertSeverityCritical},
	"on_hold":          {"Orden en espera", "qued\u00f3 en espera y necesita revisi\u00f3n", entities.AlertSeverityWarning},
	"refunded":         {"Orden reembolsada", "fue reembolsada", entities.AlertSeverityWarning},
	"delivery_novelty": {"Novedad de entrega", "tiene una novedad de entrega", entities.AlertSeverityWarning},
	"delivery_failed":  {"Entrega fallida", "no se pudo entregar", entities.AlertSeverityCritical},
	"returned":         {"Orden devuelta", "fue devuelta", entities.AlertSeverityWarning},
	"inventory_issue":  {"Novedad de inventario", "tiene una novedad de inventario", entities.AlertSeverityWarning},
	"delivered":        {"Orden entregada", "fue entregada", entities.AlertSeverityInfo},
}

var trackingAlertStatuses = map[string]string{
	"on_hold":  "tiene una novedad",
	"returned": "fue devuelta",
	"failed":   "falló en la entrega",
}

func buildAlert(event dtos.AlertEvent) (*entities.Alert, bool) {
	data := event.Data
	if data == nil {
		data = map[string]any{}
	}
	orderNumber := firstString(data, "order_number", "internal_number", "external_id")
	orderLabel := "una orden"
	if orderNumber != "" {
		orderLabel = "la orden " + orderNumber
	}
	customer := firstString(data, "customer_name")
	alert := entities.Alert{
		EventID:    event.ID,
		EventType:  event.Type,
		BusinessID: event.BusinessID,
		Severity:   entities.AlertSeverityWarning,
		CreatedAt:  event.Timestamp,
	}

	switch event.Type {
	case "order.cancelled":
		alert.Title = "Orden cancelada"
		alert.Body = fmt.Sprintf("Se canceló %s", orderLabel)
		if customer != "" {
			alert.Body += " de " + customer
		}
		if reason := firstString(data, "cancellation_reason", "reason"); reason != "" {
			alert.Body += ". Motivo: " + reason
		}
		if source := firstString(data, "source"); source == "whatsapp" {
			alert.Title = "El cliente pidió cancelar"
			alert.Body = fmt.Sprintf("%s pidió cancelar %s por WhatsApp", nameOr(customer, "El cliente"), orderLabel)
		}
		alert.Body += "."
		withOrder(&alert, orderNumber)

	case "order.status_changed":
		status := firstString(data, "current_status")
		phrase, ok := orderStatusPhrases[status]
		if !ok {
			return nil, false
		}
		alert.Title = phrase.title
		alert.Body = fmt.Sprintf("%s %s", capitalize(orderLabel), phrase.body)
		if customer != "" {
			alert.Body += " (" + customer + ")"
		}
		alert.Body += "."
		alert.Severity = phrase.severity
		withOrder(&alert, orderNumber)

	case "shipment.guide_failed":
		alert.Title = "No se pudo generar la guía"
		alert.Body = "La transportadora rechazó la guía"
		if orderNumber != "" {
			alert.Body += " de " + orderLabel
		}
		if msg := firstString(data, "error_message"); msg != "" {
			alert.Body += ": " + truncateRunes(msg, 160)
		}
		alert.Body += "."
		alert.Severity = entities.AlertSeverityCritical
		withShipment(&alert, firstString(data, "tracking_number"), orderNumber)

	case "shipment.cancel_failed":
		alert.Title = "No se pudo cancelar la guía"
		alert.Body = "La transportadora no aceptó la cancelación de la guía"
		if msg := firstString(data, "error_message"); msg != "" {
			alert.Body += ": " + truncateRunes(msg, 160)
		}
		alert.Body += "."
		withShipment(&alert, firstString(data, "tracking_number"), orderNumber)

	case "shipment.tracking_updated":
		tracking, _ := data["tracking"].(map[string]any)
		if tracking == nil {
			tracking = data
		}
		status := firstString(tracking, "probability_status", "new_status", "status")
		phrase, ok := trackingAlertStatuses[status]
		if !ok && !boolValue(tracking, "has_incidence") {
			return nil, false
		}
		if phrase == "" {
			phrase = "reportó una novedad"
		}
		guide := firstString(tracking, "tracking_number")
		orderNumber = firstString(tracking, "order_number")
		if orderNumber != "" {
			orderLabel = "la orden " + orderNumber
		}
		alert.Title = "Novedad en un envío"
		alert.Body = fmt.Sprintf("La guía %s de %s %s", nameOr(guide, "sin número"), orderLabel, phrase)
		if carrier := firstString(tracking, "carrier"); carrier != "" {
			alert.Body += " con " + carrier
		}
		if desc := firstString(tracking, "event_description", "raw_status_detail"); desc != "" {
			alert.Body += ": " + truncateRunes(desc, 160)
		}
		alert.Body += "."
		withShipment(&alert, guide, orderNumber)

	case "whatsapp.message_received":
		if firstString(data, "direction") == "outbound" {
			return nil, false
		}
		conversationID := firstString(data, "conversation_id")
		if conversationID == "" {
			return nil, false
		}
		phone := firstString(data, "phone_number")
		content := strings.TrimSpace(firstString(data, "content"))
		if content == "" {
			content = "(archivo adjunto)"
		}
		alert.Title = "Nuevo mensaje de WhatsApp"
		alert.Body = whatsAppAlertBody("", phone, content)
		alert.Severity = entities.AlertSeverityInfo
		alert.DestinationKey = "notifications"
		alert.DestinationRoute = "/notification-config?tab=conversations&conversation=" + conversationID
		alert.ReferenceType = "conversation"
		alert.ReferenceID = conversationID

	case "invoice.failed":
		alert.Title = "Factura rechazada"
		alert.Body = "No se pudo emitir la factura electrónica"
		if orderNumber != "" {
			alert.Body += " de " + orderLabel
		}
		if msg := firstString(data, "error_message", "error", "message"); msg != "" {
			alert.Body += ": " + truncateRunes(msg, 160)
		}
		alert.Body += "."
		alert.DestinationKey = "invoicing"
		alert.ReferenceType = "order"
		alert.ReferenceID = orderNumber

	case "wallet.low_balance":
		alert.Title = "Saldo bajo en la billetera"
		alert.Body = "El saldo de la billetera está por debajo del mínimo"
		if balance, ok := numberValue(data, "balance"); ok {
			alert.Body += fmt.Sprintf(": quedan %s", formatPesos(balance))
		}
		alert.Body += ". Recarga para seguir generando guías."
		alert.Severity = entities.AlertSeverityCritical
		alert.DestinationKey = "wallet"

	case "wallet.recharge.failed":
		alert.Title = "Recarga fallida"
		alert.Body = "Una recarga de la billetera no se completó"
		if msg := firstString(data, "error_message", "reason"); msg != "" {
			alert.Body += ": " + truncateRunes(msg, 160)
		}
		alert.Body += "."
		alert.DestinationKey = "wallet"

	default:
		return nil, false
	}
	return &alert, true
}

func withOrder(alert *entities.Alert, orderNumber string) {
	alert.DestinationKey = "orders"
	alert.ReferenceType = "order"
	alert.ReferenceID = orderNumber
}

func withShipment(alert *entities.Alert, guide, orderNumber string) {
	alert.DestinationKey = "shipments"
	if guide != "" {
		alert.ReferenceType = "shipment"
		alert.ReferenceID = guide
		return
	}
	alert.ReferenceType = "order"
	alert.ReferenceID = orderNumber
}

func firstString(data map[string]any, keys ...string) string {
	for _, key := range keys {
		switch v := data[key].(type) {
		case string:
			if s := strings.TrimSpace(v); s != "" {
				return s
			}
		case float64:
			if v != 0 {
				return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintf("%.2f", v), "00"), ".")
			}
		}
	}
	return ""
}

func boolValue(data map[string]any, key string) bool {
	v, _ := data[key].(bool)
	return v
}

func numberValue(data map[string]any, key string) (float64, bool) {
	switch v := data[key].(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	}
	return 0, false
}

func nameOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func formatPesos(amount float64) string {
	whole := fmt.Sprintf("%.0f", amount)
	var b strings.Builder
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 && r != '-' {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	return "$" + b.String()
}

func formatPhone(raw string) string {
	digits := strings.TrimLeft(strings.TrimSpace(raw), "+")
	if digits == "" {
		return ""
	}
	if strings.HasPrefix(digits, "57") && len(digits) == 12 {
		return "+57 " + digits[2:5] + " " + digits[5:8] + " " + digits[8:]
	}
	return "+" + digits
}

func whatsAppAlertBody(customerName, phone, content string) string {
	who := strings.TrimSpace(customerName)
	if who == "" {
		who = nameOr(formatPhone(phone), "Un cliente")
	} else if formatted := formatPhone(phone); formatted != "" {
		who += " (" + formatted + ")"
	}
	return fmt.Sprintf("%s escribi\u00f3: \u201c%s\u201d", who, truncateRunes(strings.TrimSpace(content), 180))
}
