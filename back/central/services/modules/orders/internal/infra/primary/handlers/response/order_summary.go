package response

import "time"

// OrderSummary representa un resumen de orden para listados HTTP
// ✅ DTO HTTP - CON TAGS (json)
type OrderSummary struct {
	ID                     string                 `json:"id"`
	CreatedAt              time.Time              `json:"created_at"`
	BusinessID             uint                   `json:"business_id"`
	IntegrationID          uint                   `json:"integration_id"`
	IntegrationType        string                 `json:"integration_type"`
	IntegrationLogoURL     *string                `json:"integration_logo_url,omitempty"`
	Platform               string                 `json:"platform"`
	ExternalID             string                 `json:"external_id"`
	OrderNumber            string                 `json:"order_number"`
	TotalAmount            float64                `json:"total_amount"`
	Currency               string                 `json:"currency"`
	TotalAmountPresentment float64                `json:"total_amount_presentment,omitempty"`
	CurrencyPresentment    string                 `json:"currency_presentment,omitempty"`
	CustomerName           string                 `json:"customer_name"`
	CustomerEmail          string                 `json:"customer_email"`
	CustomerPhone          string                 `json:"customer_phone,omitempty"`
	ShippingStreet         string                 `json:"shipping_street,omitempty"`
	ShippingCity           string                 `json:"shipping_city,omitempty"`
	ShippingState          string                 `json:"shipping_state,omitempty"`
	Weight                 *float64               `json:"weight,omitempty"`
	Height                 *float64               `json:"height,omitempty"`
	Width                  *float64               `json:"width,omitempty"`
	Length                 *float64               `json:"length,omitempty"`
	Status                 string                 `json:"status"`
	ItemsCount             int                    `json:"items_count"`
	DeliveryProbability    *float64               `json:"delivery_probability"`
	NegativeFactors        []string               `json:"negative_factors"`
	OrderStatus            *OrderStatusInfo       `json:"order_status,omitempty"`
	PaymentStatus          *PaymentStatusInfo     `json:"payment_status,omitempty"`
	FulfillmentStatus      *FulfillmentStatusInfo `json:"fulfillment_status,omitempty"`
	IsConfirmed            *bool                  `json:"is_confirmed"`
	Novelty                *string                `json:"novelty"`
}
