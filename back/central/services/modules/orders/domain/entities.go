package domain

import (
	"time"

	"gorm.io/datatypes"
)

// ───────────────────────────────────────────
//
//	ORDER DTOs
//
// ───────────────────────────────────────────

// CreateOrderRequest representa la solicitud para crear una orden
type CreateOrderRequest struct {
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

	// Datos estructurados (JSONB)
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

// UpdateOrderRequest representa la solicitud para actualizar una orden
type UpdateOrderRequest struct {
	// Información financiera
	Subtotal     *float64 `json:"subtotal" binding:"omitempty,min=0"`
	Tax          *float64 `json:"tax" binding:"omitempty,min=0"`
	Discount     *float64 `json:"discount" binding:"omitempty,min=0"`
	ShippingCost *float64 `json:"shipping_cost" binding:"omitempty,min=0"`
	TotalAmount  *float64 `json:"total_amount" binding:"omitempty,min=0"`
	Currency     *string  `json:"currency" binding:"omitempty,max=10"`
	CodTotal     *float64 `json:"cod_total"`

	// Información del cliente
	CustomerName  *string `json:"customer_name" binding:"omitempty,max=255"`
	CustomerEmail *string `json:"customer_email" binding:"omitempty,max=255"`
	CustomerPhone *string `json:"customer_phone" binding:"omitempty,max=32"`
	CustomerDNI   *string `json:"customer_dni" binding:"omitempty,max=64"`

	// Dirección de envío
	ShippingStreet     *string  `json:"shipping_street" binding:"omitempty,max=255"`
	ShippingCity       *string  `json:"shipping_city" binding:"omitempty,max=128"`
	ShippingState      *string  `json:"shipping_state" binding:"omitempty,max=128"`
	ShippingCountry    *string  `json:"shipping_country" binding:"omitempty,max=128"`
	ShippingPostalCode *string  `json:"shipping_postal_code" binding:"omitempty,max=32"`
	ShippingLat        *float64 `json:"shipping_lat"`
	ShippingLng        *float64 `json:"shipping_lng"`

	// Información de pago
	PaymentMethodID *uint      `json:"payment_method_id"`
	IsPaid          *bool      `json:"is_paid"`
	PaidAt          *time.Time `json:"paid_at"`

	// Información de envío/logística
	TrackingNumber *string    `json:"tracking_number"`
	TrackingLink   *string    `json:"tracking_link"`
	GuideID        *string    `json:"guide_id"`
	GuideLink      *string    `json:"guide_link"`
	DeliveryDate   *time.Time `json:"delivery_date"`
	DeliveredAt    *time.Time `json:"delivered_at"`

	// Información de fulfillment
	WarehouseID   *uint   `json:"warehouse_id"`
	WarehouseName *string `json:"warehouse_name" binding:"omitempty,max=128"`
	DriverID      *uint   `json:"driver_id"`
	DriverName    *string `json:"driver_name" binding:"omitempty,max=255"`
	IsLastMile    *bool   `json:"is_last_mile"`

	// Dimensiones y peso
	Weight *float64 `json:"weight"`
	Height *float64 `json:"height"`
	Width  *float64 `json:"width"`
	Length *float64 `json:"length"`
	Boxes  *string  `json:"boxes"`

	// Tipo y estado
	OrderTypeID    *uint   `json:"order_type_id"`
	OrderTypeName  *string `json:"order_type_name" binding:"omitempty,max=64"`
	Status         *string `json:"status" binding:"omitempty,max=64"`
	OriginalStatus *string `json:"original_status" binding:"omitempty,max=64"`

	// Información adicional
	Notes    *string `json:"notes"`
	Coupon   *string `json:"coupon"`
	Approved *bool   `json:"approved"`
	UserID   *uint   `json:"user_id"`
	UserName *string `json:"user_name" binding:"omitempty,max=255"`

	// Facturación
	Invoiceable     *bool   `json:"invoiceable"`
	InvoiceURL      *string `json:"invoice_url"`
	InvoiceID       *string `json:"invoice_id"`
	InvoiceProvider *string `json:"invoice_provider"`

	// Datos estructurados (JSONB)
	Items              datatypes.JSON `json:"items"`
	Metadata           datatypes.JSON `json:"metadata"`
	FinancialDetails   datatypes.JSON `json:"financial_details"`
	ShippingDetails    datatypes.JSON `json:"shipping_details"`
	PaymentDetails     datatypes.JSON `json:"payment_details"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details"`
}

// OrderResponse representa la respuesta de una orden
type OrderResponse struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

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
	CustomerID    *uint  `json:"customer_id,omitempty"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	CustomerDNI   string `json:"customer_dni"`

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

	// Datos estructurados (JSONB)
	Items              datatypes.JSON `json:"items,omitempty"`
	Metadata           datatypes.JSON `json:"metadata,omitempty"`
	FinancialDetails   datatypes.JSON `json:"financial_details,omitempty"`
	ShippingDetails    datatypes.JSON `json:"shipping_details,omitempty"`
	PaymentDetails     datatypes.JSON `json:"payment_details,omitempty"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details,omitempty"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`
}

// OrdersListResponse representa la respuesta paginada de órdenes
type OrdersListResponse struct {
	Data       []OrderResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
