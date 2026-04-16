package response

import "time"

type OrderEventMessage struct {
	EventID       string         `json:"event_id"`
	EventType     string         `json:"event_type"`
	OrderID       string         `json:"order_id"`
	BusinessID    *uint          `json:"business_id"`
	IntegrationID *uint          `json:"integration_id"`
	Timestamp     time.Time      `json:"timestamp"`
	Order         *OrderSnapshot `json:"order"`
	Changes       map[string]any `json:"changes,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

type OrderSnapshot struct {
	ID             string `json:"id"`
	OrderNumber    string `json:"order_number"`
	InternalNumber string `json:"internal_number"`
	ExternalID     string `json:"external_id"`

	TotalAmount     float64 `json:"total_amount"`
	Currency        string  `json:"currency"`
	PaymentMethodID uint    `json:"payment_method_id"`
	PaymentStatusID *uint   `json:"payment_status_id,omitempty"`

	Subtotal     float64 `json:"subtotal"`
	Tax          float64 `json:"tax"`
	Discount     float64 `json:"discount"`
	ShippingCost float64 `json:"shipping_cost"`

	CustomerID    *uint  `json:"customer_id,omitempty"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email,omitempty"`
	CustomerPhone string `json:"customer_phone,omitempty"`
	CustomerDNI   string `json:"customer_dni,omitempty"`

	Platform      string `json:"platform"`
	IntegrationID uint   `json:"integration_id"`
	BusinessName  string `json:"business_name,omitempty"`

	OrderStatusID       *uint `json:"order_status_id,omitempty"`
	FulfillmentStatusID *uint `json:"fulfillment_status_id,omitempty"`

	WarehouseID *uint `json:"warehouse_id,omitempty"`

	Items []OrderItemSnapshot `json:"items,omitempty"`

	ShippingStreet      string   `json:"shipping_street,omitempty"`
	ShippingCity        string   `json:"shipping_city,omitempty"`
	ShippingState       string   `json:"shipping_state,omitempty"`
	ShippingCountry     string   `json:"shipping_country,omitempty"`
	ShippingPostalCode  string   `json:"shipping_postal_code,omitempty"`
	ShippingLat         *float64 `json:"shipping_lat,omitempty"`
	ShippingLng         *float64 `json:"shipping_lng,omitempty"`
	ItemsSummary        string  `json:"items_summary,omitempty"`
	ShippingAddress     string  `json:"shipping_address,omitempty"`
	IsPaid              bool    `json:"is_paid"`
	DeliveryProbability float64 `json:"delivery_probability,omitempty"`
	Status              string  `json:"status,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderItemSnapshot struct {
	ProductID *string `json:"product_id,omitempty"`
	SKU       string  `json:"sku"`
	VariantID *string `json:"variant_id,omitempty"`

	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`

	Tax             float64  `json:"tax"`
	TaxRate         *float64 `json:"tax_rate,omitempty"`
	Discount        float64  `json:"discount"`
	DiscountPercent float64  `json:"discount_percent"`

	ImageURL   *string `json:"image_url,omitempty"`
	ProductURL *string `json:"product_url,omitempty"`
}
