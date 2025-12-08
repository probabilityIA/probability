package domain

import (
	"crypto/rand"
	"time"
)

// ───────────────────────────────────────────
//
//	ORDER EVENT TYPES
//
// ───────────────────────────────────────────

// OrderEventType define los tipos de eventos relacionados con órdenes
type OrderEventType string

const (
	// Eventos de ciclo de vida de la orden
	OrderEventTypeCreated         OrderEventType = "order.created"
	OrderEventTypeUpdated         OrderEventType = "order.updated"
	OrderEventTypeStatusChanged   OrderEventType = "order.status_changed"
	OrderEventTypeCancelled       OrderEventType = "order.cancelled"
	OrderEventTypeDelivered       OrderEventType = "order.delivered"
	OrderEventTypeShipped         OrderEventType = "order.shipped"
	OrderEventTypePaymentReceived OrderEventType = "order.payment_received"
	OrderEventTypeRefunded        OrderEventType = "order.refunded"
	OrderEventTypeFailed          OrderEventType = "order.failed"
	OrderEventTypeOnHold          OrderEventType = "order.on_hold"
	OrderEventTypeProcessing      OrderEventType = "order.processing"
)

// ───────────────────────────────────────────
//
//	ORDER EVENT STRUCTURES
//
// ───────────────────────────────────────────

// OrderEvent representa un evento relacionado con una orden
type OrderEvent struct {
	ID            string                 `json:"id"`
	Type          OrderEventType         `json:"type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id,omitempty"`
	IntegrationID *uint                  `json:"integration_id,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          OrderEventData         `json:"data"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// OrderEventData contiene los datos específicos del evento de orden
type OrderEventData struct {
	// Información básica de la orden
	OrderNumber    string `json:"order_number,omitempty"`
	InternalNumber string `json:"internal_number,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`

	// Cambios de estado
	PreviousStatus string `json:"previous_status,omitempty"`
	CurrentStatus  string `json:"current_status,omitempty"`

	// Información adicional
	CustomerEmail string                 `json:"customer_email,omitempty"`
	TotalAmount   *float64               `json:"total_amount,omitempty"`
	Currency      string                 `json:"currency,omitempty"`
	Platform      string                 `json:"platform,omitempty"`
	Extra         map[string]interface{} `json:"extra,omitempty"`
}

// ───────────────────────────────────────────
//
//	HELPER FUNCTIONS
//
// ───────────────────────────────────────────

// NewOrderEvent crea un nuevo evento de orden
func NewOrderEvent(eventType OrderEventType, orderID string, data OrderEventData) *OrderEvent {
	return &OrderEvent{
		ID:        generateEventID(),
		Type:      eventType,
		OrderID:   orderID,
		Timestamp: time.Now(),
		Data:      data,
		Metadata:  make(map[string]interface{}),
	}
}

// generateEventID genera un ID único para el evento
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString genera una cadena aleatoria
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
