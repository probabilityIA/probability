package entities

import "time"

const (
	NotificationChannelWhatsApp = "whatsapp"
	NotificationChannelEmail    = "email"
	NotificationChannelSMS      = "sms"
	NotificationChannelSSE      = "sse"
)

const (
	NotificationStateNone    = "none"
	NotificationStatePartial = "partial"
	NotificationStateDone    = "done"
)

type OrderNotificationKey struct {
	OrderID       string
	OrderNumber   string
	IntegrationID uint
}

type OrderNotificationCounter struct {
	Channel  string
	Sent     int
	Expected int
	Failed   int
	State    string
}

func NewOrderNotificationCounter(channel string, sent, expected, failed int) OrderNotificationCounter {
	if sent > expected {
		sent = expected
	}

	state := NotificationStatePartial
	switch {
	case sent == 0:
		state = NotificationStateNone
	case sent >= expected:
		state = NotificationStateDone
	}

	return OrderNotificationCounter{
		Channel:  channel,
		Sent:     sent,
		Expected: expected,
		Failed:   failed,
		State:    state,
	}
}

type OrderNotificationMessage struct {
	ID           string
	Channel      string
	TemplateName string
	Direction    string
	Status       string
	Content      string
	MessageID    string
	CreatedAt    time.Time
}

const (
	NotificationEventStatusPending = "pending"
	NotificationEventStatusSent    = "sent"
	NotificationEventStatusFailed  = "failed"
)

var WhatsAppEventTemplates = map[string][]string{
	"order.created":            {"confirmacion_pedido", "confirmacion_pedido_contraentrega"},
	"order.shipped":            {"pedido_en_reparto", "pedido_en_reparto_cod"},
	"order.delivered":          {"pedido_entregado", "pedido_entregado_cod"},
	"shipment.guide_generated": {"guia_envio_generada", "guia_envio_generada_cod"},
	"order.canceled":           {"pedido_cancelado"},
	"invoice.created":          {"factura_generada"},
}

func EventStatusFromTemplates(eventCode string, statusByTemplate map[string]string) string {
	status := NotificationEventStatusPending
	for _, tpl := range WhatsAppEventTemplates[eventCode] {
		s, ok := statusByTemplate[tpl]
		if !ok {
			continue
		}
		status = s
		if s == NotificationEventStatusSent {
			break
		}
	}
	return status
}
