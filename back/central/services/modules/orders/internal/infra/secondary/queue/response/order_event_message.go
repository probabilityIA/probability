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

	// Información financiera detallada (para facturación)
	Subtotal     float64 `json:"subtotal"`
	Tax          float64 `json:"tax"`
	Discount     float64 `json:"discount"`
	ShippingCost float64 `json:"shipping_cost"`

	// Información del cliente
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	CustomerDNI   string `json:"customer_dni,omitempty"`

	// Información de origen
	Platform      string `json:"platform"`
	IntegrationID uint   `json:"integration_id"`

	// Estados
	OrderStatusID       *uint `json:"order_status_id,omitempty"`
	FulfillmentStatusID *uint `json:"fulfillment_status_id,omitempty"`

	// Items detallados (para facturación e inventario)
	Items []OrderItemSnapshot `json:"items,omitempty"`

	// Items y envío (información adicional para mensajes)
	ItemsSummary    string `json:"items_summary,omitempty"`
	ShippingAddress string `json:"shipping_address,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderItemSnapshot representa un item de una orden en un snapshot
// Incluye toda la información necesaria para facturación, inventario y reportes
// sin necesidad de consultar la base de datos
type OrderItemSnapshot struct {
	// Identificadores
	ProductID *string `json:"product_id,omitempty"`
	SKU       string  `json:"sku"`
	VariantID *string `json:"variant_id,omitempty"`

	// Información del producto
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"` // Título alternativo (ProductTitle en BD)
	Description string `json:"description,omitempty"`

	// Cantidades y precios
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`

	// Impuestos y descuentos
	Tax      float64  `json:"tax"`
	TaxRate  *float64 `json:"tax_rate,omitempty"`
	Discount float64  `json:"discount"`

	// Información adicional
	ImageURL   *string `json:"image_url,omitempty"`
	ProductURL *string `json:"product_url,omitempty"`
}
