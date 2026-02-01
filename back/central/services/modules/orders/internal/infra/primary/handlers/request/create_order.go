package request

import (
	"time"

	"gorm.io/datatypes"
)

// CreateOrder representa la petición HTTP para crear una orden
// ✅ DTO HTTP - CON TAGS (json + binding)
type CreateOrder struct {
	// Identificadores de integración
	BusinessID      *uint  `json:"business_id"`
	IntegrationID   uint   `json:"integration_id" binding:"required"`
	IntegrationType string `json:"integration_type" binding:"required,max=50"`

	// Identificadores de la orden
	Platform       string `json:"platform" binding:"required,max=50"`
	ExternalID     string `json:"external_id" binding:"required,max=255"`
	OrderNumber    string `json:"order_number" binding:"max=128"`
	InternalNumber string `json:"internal_number" binding:"max=128"`

	// Información financiera
	Subtotal     float64  `json:"subtotal" binding:"required,min=0"`
	Tax          float64  `json:"tax" binding:"min=0"`
	Discount     float64  `json:"discount" binding:"min=0"`
	ShippingCost float64  `json:"shipping_cost" binding:"min=0"`
	TotalAmount  float64  `json:"total_amount" binding:"required,min=0"`
	Currency     string   `json:"currency" binding:"max=10"`
	CodTotal     *float64 `json:"cod_total"`

	// Información del cliente
	CustomerID    *uint  `json:"customer_id"`
	CustomerName  string `json:"customer_name" binding:"max=255"`
	CustomerEmail string `json:"customer_email" binding:"max=255"`
	CustomerPhone string `json:"customer_phone" binding:"max=32"`
	CustomerDNI   string `json:"customer_dni" binding:"max=64"`

	// Dirección de envío
	ShippingStreet     string   `json:"shipping_street" binding:"max=255"`
	ShippingCity       string   `json:"shipping_city" binding:"max=128"`
	ShippingState      string   `json:"shipping_state" binding:"max=128"`
	ShippingCountry    string   `json:"shipping_country" binding:"max=128"`
	ShippingPostalCode string   `json:"shipping_postal_code" binding:"max=32"`
	ShippingLat        *float64 `json:"shipping_lat"`
	ShippingLng        *float64 `json:"shipping_lng"`

	// Información de pago
	PaymentMethodID uint       `json:"payment_method_id" binding:"required"`
	IsPaid          bool       `json:"is_paid"`
	PaidAt          *time.Time `json:"paid_at"`

	// Información de envío/logística
	TrackingNumber *string    `json:"tracking_number"`
	TrackingLink   *string    `json:"tracking_link"`
	GuideID        *string    `json:"guide_id"`
	GuideLink      *string    `json:"guide_link"`
	DeliveryDate   *time.Time `json:"delivery_date"`
	DeliveredAt    *time.Time `json:"delivered_at"`

	// Información de fulfillment
	WarehouseID   *uint  `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name" binding:"max=128"`
	DriverID      *uint  `json:"driver_id"`
	DriverName    string `json:"driver_name" binding:"max=255"`
	IsLastMile    bool   `json:"is_last_mile"`

	// Dimensiones y peso
	Weight *float64 `json:"weight"`
	Height *float64 `json:"height"`
	Width  *float64 `json:"width"`
	Length *float64 `json:"length"`
	Boxes  *string  `json:"boxes"`

	// Tipo y estado
	OrderTypeID    *uint  `json:"order_type_id"`
	OrderTypeName  string `json:"order_type_name" binding:"max=64"`
	Status         string `json:"status" binding:"max=64"`
	OriginalStatus string `json:"original_status" binding:"max=64"`
	StatusID       *uint  `json:"status_id"`

	// Estados independientes
	PaymentStatusID     *uint `json:"payment_status_id"`
	FulfillmentStatusID *uint `json:"fulfillment_status_id"`

	// Información adicional
	Notes    *string `json:"notes"`
	Coupon   *string `json:"coupon"`
	Approved *bool   `json:"approved"`
	UserID   *uint   `json:"user_id"`
	UserName string  `json:"user_name" binding:"max=255"`

	// Facturación
	Invoiceable     bool    `json:"invoiceable"`
	InvoiceURL      *string `json:"invoice_url"`
	InvoiceID       *string `json:"invoice_id"`
	InvoiceProvider *string `json:"invoice_provider"`

	// Datos estructurados (JSONB) - permitido en infra
	Items              datatypes.JSON `json:"items"`
	Metadata           datatypes.JSON `json:"metadata"`
	FinancialDetails   datatypes.JSON `json:"financial_details"`
	ShippingDetails    datatypes.JSON `json:"shipping_details"`
	PaymentDetails     datatypes.JSON `json:"payment_details"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`
}
