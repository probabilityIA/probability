package response

import "time"

// OrderEventMessage representa el payload unificado de eventos de órdenes en RabbitMQ
type OrderEventMessage struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	OrderID       string                 `json:"order_id"`
	BusinessID    *uint                  `json:"business_id"`
	IntegrationID *uint                  `json:"integration_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Order         *OrderSnapshot         `json:"order"`      // Snapshot completo SIEMPRE incluido
	Changes       map[string]interface{} `json:"changes,omitempty"`    // Cambios específicos
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// OrderSnapshot representa un snapshot completo de una orden
// Incluye toda la información necesaria para que consumidores externos
// (WhatsApp, Invoicing, etc.) puedan decidir si actuar sin consultar BD
type OrderSnapshot struct {
	// Identificadores
	ID             string `json:"id"`
	OrderNumber    string `json:"order_number"`
	InternalNumber string `json:"internal_number"`
	ExternalID     string `json:"external_id"`

	// Información financiera
	TotalAmount     float64 `json:"total_amount"`
	Currency        string  `json:"currency"`
	PaymentMethodID uint    `json:"payment_method_id"`
	PaymentStatusID *uint   `json:"payment_status_id,omitempty"`

	// Información del cliente
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`

	// Información de origen
	Platform      string `json:"platform"`
	IntegrationID uint   `json:"integration_id"`

	// Estados
	OrderStatusID       *uint `json:"order_status_id,omitempty"`
	FulfillmentStatusID *uint `json:"fulfillment_status_id,omitempty"`

	// Items y envío (información adicional para mensajes)
	ItemsSummary    string `json:"items_summary,omitempty"`
	ShippingAddress string `json:"shipping_address,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
