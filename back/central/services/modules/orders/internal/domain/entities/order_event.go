package entities

import (
	"crypto/rand"
	"time"
)

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

	// Eventos de cálculo de score
	OrderEventTypeScoreCalculationRequested OrderEventType = "order.score_calculation_requested"
	OrderEventTypeScoreCalculated          OrderEventType = "order.score_calculated"

	// Eventos de confirmación
	OrderEventTypeConfirmationRequested OrderEventType = "order.confirmation_requested"
)

// OrderEvent representa un evento relacionado con una orden
type OrderEvent struct {
	ID            string
	Type          OrderEventType
	OrderID       string
	BusinessID    *uint
	IntegrationID *uint
	Timestamp     time.Time
	Data          OrderEventData
	Metadata      map[string]interface{}
}

// OrderEventData contiene los datos específicos del evento de orden
type OrderEventData struct {
	// Información básica de la orden
	OrderNumber    string
	InternalNumber string
	ExternalID     string

	// Cambios de estado
	PreviousStatus string
	CurrentStatus  string

	// Información adicional
	CustomerEmail string
	TotalAmount   *float64
	Currency      string
	Platform      string
	Extra         map[string]interface{}
}

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
