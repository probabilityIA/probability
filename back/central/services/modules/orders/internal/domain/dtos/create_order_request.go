package dtos

import (
	"encoding/json"
	"time"
)

// CreateOrderRequest representa la solicitud para crear una orden
// ✅ DTO HTTP/DOMAIN - CON TAGS
type CreateOrderRequest struct {
	// Identificadores de integración
	BusinessID      *uint  `json:"business_id"`
	IntegrationID   uint   `json:"integration_id"`
	IntegrationType string `json:"integration_type"`

	// Identificadores de la orden
	Platform       string `json:"platform"`
	ExternalID     string `json:"external_id"`
	OrderNumber    string `json:"order_number"`
	InternalNumber string `json:"internal_number"`

	// Información financiera
	Subtotal     float64  `json:"subtotal"`
	Tax          float64  `json:"tax"`
	Discount     float64  `json:"discount"`
	ShippingCost float64  `json:"shipping_cost"`
	TotalAmount  float64  `json:"total_amount"`
	Currency     string   `json:"currency"`
	CodTotal     *float64 `json:"cod_total,omitempty"`

	// Información del cliente
	CustomerID        *uint  `json:"customer_id,omitempty"`
	CustomerName      string `json:"customer_name"`
	CustomerFirstName string `json:"customer_first_name"`
	CustomerLastName  string `json:"customer_last_name"`
	CustomerEmail     string `json:"customer_email"`
	CustomerPhone     string `json:"customer_phone"`
	CustomerDNI       string `json:"customer_dni"`

	// Dirección de envío
	ShippingStreet     string   `json:"shipping_street"`
	ShippingCity       string   `json:"shipping_city"`
	ShippingState      string   `json:"shipping_state"`
	ShippingCountry    string   `json:"shipping_country"`
	ShippingPostalCode string   `json:"shipping_postal_code"`
	ShippingLat        *float64 `json:"shipping_lat,omitempty"`
	ShippingLng        *float64 `json:"shipping_lng,omitempty"`

	// Información de pago
	PaymentMethodID uint       `json:"payment_method_id"`
	IsPaid          bool       `json:"is_paid"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`

	// Información de envío/logística
	TrackingNumber *string    `json:"tracking_number,omitempty"`
	TrackingLink   *string    `json:"tracking_link,omitempty"`
	GuideID        *string    `json:"guide_id,omitempty"`
	GuideLink      *string    `json:"guide_link,omitempty"`
	DeliveryDate   *time.Time `json:"delivery_date,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`

	// Información de fulfillment
	WarehouseID   *uint  `json:"warehouse_id,omitempty"`
	WarehouseName string `json:"warehouse_name"`
	DriverID      *uint  `json:"driver_id,omitempty"`
	DriverName    string `json:"driver_name"`
	IsLastMile    bool   `json:"is_last_mile"`

	// Dimensiones y peso
	Weight *float64 `json:"weight,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Width  *float64 `json:"width,omitempty"`
	Length *float64 `json:"length,omitempty"`
	Boxes  *string  `json:"boxes,omitempty"`

	// Tipo y estado
	OrderTypeID    *uint  `json:"order_type_id,omitempty"`
	OrderTypeName  string `json:"order_type_name"`
	Status         string `json:"status"`
	OriginalStatus string `json:"original_status"`
	StatusID       *uint  `json:"status_id,omitempty"`

	// Estados independientes
	PaymentStatusID     *uint `json:"payment_status_id,omitempty"`
	FulfillmentStatusID *uint `json:"fulfillment_status_id,omitempty"`

	// Información adicional
	Notes    *string `json:"notes,omitempty"`
	Coupon   *string `json:"coupon,omitempty"`
	Approved *bool   `json:"approved,omitempty"`
	UserID   *uint   `json:"user_id,omitempty"`
	UserName string  `json:"user_name"`

	// Facturación
	Invoiceable     bool    `json:"invoiceable"`
	InvoiceURL      *string `json:"invoice_url,omitempty"`
	InvoiceID       *string `json:"invoice_id,omitempty"`
	InvoiceProvider *string `json:"invoice_provider,omitempty"`

	// Datos estructurados (JSON) - json.RawMessage para preservar el formato
	Items              json.RawMessage `json:"items,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	FinancialDetails   json.RawMessage `json:"financial_details,omitempty"`
	ShippingDetails    json.RawMessage `json:"shipping_details,omitempty"`
	PaymentDetails     json.RawMessage `json:"payment_details,omitempty"`
	FulfillmentDetails json.RawMessage `json:"fulfillment_details,omitempty"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`
}
