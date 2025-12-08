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

// OrderSummary representa un resumen de la orden para listados
type OrderSummary struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	BusinessID      uint      `json:"business_id"`
	IntegrationID   uint      `json:"integration_id"`
	IntegrationType string    `json:"integration_type"`
	Platform        string    `json:"platform"`
	ExternalID      string    `json:"external_id"`
	OrderNumber     string    `json:"order_number"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	Status          string    `json:"status"`
	PaymentStatus   string    `json:"payment_status"` // derived from IsPaid
	ItemsCount      int       `json:"items_count"`    // derived from len(Items)
}

// OrderRawResponse representa la respuesta con los datos crudos
type OrderRawResponse struct {
	OrderID       string         `json:"order_id"`
	ChannelSource string         `json:"channel_source"`
	RawData       datatypes.JSON `json:"raw_data"`
}

// OrdersListResponse representa la respuesta paginada de órdenes
type OrdersListResponse struct {
	Data       []OrderSummary `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// ───────────────────────────────────────────
//
//	CANONICAL ORDER DTO - Estructura canónica para recibir órdenes mapeadas
//
// ───────────────────────────────────────────

// CanonicalOrderDTO representa la estructura canónica que todas las integraciones
// deben enviar después de mapear sus datos específicos
type CanonicalOrderDTO struct {
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

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`

	// Datos estructurados (JSONB) - Para compatibilidad
	Items              datatypes.JSON `json:"items,omitempty"`
	Metadata           datatypes.JSON `json:"metadata,omitempty"`
	FinancialDetails   datatypes.JSON `json:"financial_details,omitempty"`
	ShippingDetails    datatypes.JSON `json:"shipping_details,omitempty"`
	PaymentDetails     datatypes.JSON `json:"payment_details,omitempty"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details,omitempty"`

	// ============================================
	// TABLAS RELACIONADAS
	// ============================================

	// Items de la orden
	OrderItems []CanonicalOrderItemDTO `json:"order_items" binding:"dive"`

	// Direcciones
	Addresses []CanonicalAddressDTO `json:"addresses" binding:"dive"`

	// Pagos
	Payments []CanonicalPaymentDTO `json:"payments" binding:"dive"`

	// Envíos
	Shipments []CanonicalShipmentDTO `json:"shipments" binding:"dive"`

	// Metadata del canal (datos crudos)
	ChannelMetadata *CanonicalChannelMetadataDTO `json:"channel_metadata"`
}

// CanonicalOrderItemDTO representa un item/producto de la orden
type CanonicalOrderItemDTO struct {
	ProductID    *string        `json:"product_id"`
	ProductSKU   string         `json:"product_sku" binding:"required,max=128"`
	ProductName  string         `json:"product_name" binding:"required,max=255"`
	ProductTitle string         `json:"product_title" binding:"max=255"`
	VariantID    *string        `json:"variant_id"`
	Quantity     int            `json:"quantity" binding:"required,min=1"`
	UnitPrice    float64        `json:"unit_price" binding:"required,min=0"`
	TotalPrice   float64        `json:"total_price" binding:"required,min=0"`
	Currency     string         `json:"currency" binding:"max=10"`
	Discount     float64        `json:"discount" binding:"min=0"`
	Tax          float64        `json:"tax" binding:"min=0"`
	TaxRate      *float64       `json:"tax_rate"`
	ImageURL     *string        `json:"image_url"`
	ProductURL   *string        `json:"product_url"`
	Weight       *float64       `json:"weight"`
	Metadata     datatypes.JSON `json:"metadata,omitempty"`
}

// CanonicalAddressDTO representa una dirección (envío o facturación)
type CanonicalAddressDTO struct {
	Type         string         `json:"type" binding:"required,oneof=shipping billing"` // "shipping" o "billing"
	FirstName    string         `json:"first_name" binding:"max=128"`
	LastName     string         `json:"last_name" binding:"max=128"`
	Company      string         `json:"company" binding:"max=255"`
	Phone        string         `json:"phone" binding:"max=32"`
	Street       string         `json:"street" binding:"required,max=255"`
	Street2      string         `json:"street2" binding:"max=255"`
	City         string         `json:"city" binding:"required,max=128"`
	State        string         `json:"state" binding:"max=128"`
	Country      string         `json:"country" binding:"required,max=128"`
	PostalCode   string         `json:"postal_code" binding:"max=32"`
	Latitude     *float64       `json:"latitude"`
	Longitude    *float64       `json:"longitude"`
	Instructions *string        `json:"instructions"`
	Metadata     datatypes.JSON `json:"metadata,omitempty"`
}

// CanonicalPaymentDTO representa un pago de la orden
type CanonicalPaymentDTO struct {
	PaymentMethodID  uint           `json:"payment_method_id" binding:"required"`
	Amount           float64        `json:"amount" binding:"required,min=0"`
	Currency         string         `json:"currency" binding:"max=10"`
	ExchangeRate     *float64       `json:"exchange_rate"`
	Status           string         `json:"status" binding:"required,oneof=pending completed failed refunded"`
	PaidAt           *time.Time     `json:"paid_at"`
	ProcessedAt      *time.Time     `json:"processed_at"`
	TransactionID    *string        `json:"transaction_id"`
	PaymentReference *string        `json:"payment_reference"`
	Gateway          *string        `json:"gateway"`
	RefundAmount     *float64       `json:"refund_amount"`
	RefundedAt       *time.Time     `json:"refunded_at"`
	FailureReason    *string        `json:"failure_reason"`
	Metadata         datatypes.JSON `json:"metadata,omitempty"`
}

// CanonicalShipmentDTO representa un envío de la orden
type CanonicalShipmentDTO struct {
	TrackingNumber    *string        `json:"tracking_number"`
	TrackingURL       *string        `json:"tracking_url"`
	Carrier           *string        `json:"carrier"`
	CarrierCode       *string        `json:"carrier_code"`
	GuideID           *string        `json:"guide_id"`
	GuideURL          *string        `json:"guide_url"`
	Status            string         `json:"status" binding:"oneof=pending in_transit delivered failed"`
	ShippedAt         *time.Time     `json:"shipped_at"`
	DeliveredAt       *time.Time     `json:"delivered_at"`
	ShippingAddressID *uint          `json:"shipping_address_id"`
	ShippingCost      *float64       `json:"shipping_cost"`
	InsuranceCost     *float64       `json:"insurance_cost"`
	TotalCost         *float64       `json:"total_cost"`
	Weight            *float64       `json:"weight"`
	Height            *float64       `json:"height"`
	Width             *float64       `json:"width"`
	Length            *float64       `json:"length"`
	WarehouseID       *uint          `json:"warehouse_id"`
	WarehouseName     string         `json:"warehouse_name" binding:"max=128"`
	DriverID          *uint          `json:"driver_id"`
	DriverName        string         `json:"driver_name" binding:"max=255"`
	IsLastMile        bool           `json:"is_last_mile"`
	EstimatedDelivery *time.Time     `json:"estimated_delivery"`
	DeliveryNotes     *string        `json:"delivery_notes"`
	Metadata          datatypes.JSON `json:"metadata,omitempty"`
}

// CanonicalChannelMetadataDTO representa los datos crudos del canal
type CanonicalChannelMetadataDTO struct {
	ChannelSource string         `json:"channel_source" binding:"required,max=50"`
	RawData       datatypes.JSON `json:"raw_data" binding:"required"`
	Version       string         `json:"version" binding:"max=20"`
	ReceivedAt    time.Time      `json:"received_at"`
	ProcessedAt   *time.Time     `json:"processed_at"`
	IsLatest      bool           `json:"is_latest"`
	LastSyncedAt  *time.Time     `json:"last_synced_at"`
	SyncStatus    string         `json:"sync_status" binding:"max=64"`
}
