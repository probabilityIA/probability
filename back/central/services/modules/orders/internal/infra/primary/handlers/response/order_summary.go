package response

import (
	"encoding/json"
	"time"
)

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
	ChannelPackID          string                 `json:"channel_pack_id,omitempty"`
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
	ShippingGeoConfidence  string                 `json:"shipping_geo_confidence,omitempty"`
	Weight                 *float64               `json:"weight,omitempty"`
	Height                 *float64               `json:"height,omitempty"`
	Width                  *float64               `json:"width,omitempty"`
	Length                 *float64               `json:"length,omitempty"`
	Status                 string                 `json:"status"`
	ItemsCount             int                    `json:"items_count"`
	DeliveryProbability    *float64               `json:"delivery_probability"`
	NegativeFactors        []string               `json:"negative_factors"`
	ScoreBreakdown         json.RawMessage        `json:"score_breakdown,omitempty"`
	OrderStatus            *OrderStatusInfo       `json:"order_status,omitempty"`
	PaymentStatus          *PaymentStatusInfo     `json:"payment_status,omitempty"`
	FulfillmentStatus      *FulfillmentStatusInfo `json:"fulfillment_status,omitempty"`
	OrderStatusURL         string                 `json:"order_status_url,omitempty"`
	GuideLink              *string                `json:"guide_link,omitempty"`
	IsPaid                 bool                   `json:"is_paid"`
	IsCod                  bool                   `json:"is_cod"`
	CodTotal               *float64               `json:"cod_total,omitempty"`
	IsConfirmed            *bool                  `json:"is_confirmed"`
	Novelty                *string                `json:"novelty"`
	IsTest                 bool                   `json:"is_test"`
	InvoiceStatus          string                 `json:"invoice_status"`
	CodCutConfirmed        bool                   `json:"cod_cut_confirmed"`
	Shipment               *ShipmentSummary       `json:"shipment,omitempty"`
	QuotedShipping         *QuotedShipping        `json:"quoted_shipping,omitempty"`
	FreeShipping           bool                   `json:"free_shipping"`
	StatusSource           string                 `json:"status_source,omitempty"`
	StatusChangedBy        string                 `json:"status_changed_by,omitempty"`
	StatusChangedAt        *time.Time             `json:"status_changed_at,omitempty"`
	Notifications          []NotificationCounter  `json:"notifications,omitempty"`
}

type NotificationCounter struct {
	Channel  string `json:"channel"`
	Sent     int    `json:"sent"`
	Expected int    `json:"expected"`
	Failed   int    `json:"failed"`
	State    string `json:"state"`
}
