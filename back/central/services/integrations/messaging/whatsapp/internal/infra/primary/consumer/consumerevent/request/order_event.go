package request

import "time"

// OrderEventType representa el tipo de evento de orden
type OrderEventType string

// OrderEvent representa un evento de orden desde RabbitMQ
type OrderEvent struct {
	ID            string                 `json:"id"`
	Type          OrderEventType         `json:"event_type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id,omitempty"`
	IntegrationID *uint                  `json:"integration_id,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          OrderEventData         `json:"data"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// OrderEventData contiene los datos específicos del evento de orden
type OrderEventData struct {
	OrderNumber    string   `json:"order_number,omitempty"`
	InternalNumber string   `json:"internal_number,omitempty"`
	ExternalID     string   `json:"external_id,omitempty"`
	PreviousStatus string   `json:"previous_status,omitempty"`
	CurrentStatus  string   `json:"current_status,omitempty"`
	CustomerEmail  string   `json:"customer_email,omitempty"`
	TotalAmount    *float64 `json:"total_amount,omitempty"`
	Currency       string   `json:"currency,omitempty"`
	Platform       string   `json:"platform,omitempty"`
}
