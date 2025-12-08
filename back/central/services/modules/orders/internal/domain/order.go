package domain

import (
	"time"

	"gorm.io/datatypes"
)

// Order representa una orden en el dominio
type Order struct {
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
	CodTotal     *float64 `json:"cod_total"`

	// Información del cliente
	CustomerID    *uint  `json:"customer_id"`
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
	ShippingLat        *float64 `json:"shipping_lat"`
	ShippingLng        *float64 `json:"shipping_lng"`

	// Información de pago
	PaymentMethodID uint       `json:"payment_method_id"`
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
	WarehouseName string `json:"warehouse_name"`
	DriverID      *uint  `json:"driver_id"`
	DriverName    string `json:"driver_name"`
	IsLastMile    bool   `json:"is_last_mile"`

	// Dimensiones y peso
	Weight *float64 `json:"weight"`
	Height *float64 `json:"height"`
	Width  *float64 `json:"width"`
	Length *float64 `json:"length"`
	Boxes  *string  `json:"boxes"`

	// Tipo y estado
	OrderTypeID    *uint  `json:"order_type_id"`
	OrderTypeName  string `json:"order_type_name"`
	Status         string `json:"status"`
	OriginalStatus string `json:"original_status"`

	// Información adicional
	Notes    *string `json:"notes"`
	Coupon   *string `json:"coupon"`
	Approved *bool   `json:"approved"`
	UserID   *uint   `json:"user_id"`
	UserName string  `json:"user_name"`

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

	// Relaciones
	OrderItems      []OrderItem            `json:"order_items"`
	Addresses       []Address              `json:"addresses"`
	Payments        []Payment              `json:"payments"`
	Shipments       []Shipment             `json:"shipments"`
	ChannelMetadata []OrderChannelMetadata `json:"channel_metadata"`
}

// OrderItem representa un item de la orden en el dominio
type OrderItem struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID string `json:"order_id"`

	ProductID    *string `json:"product_id"`
	ProductSKU   string  `json:"product_sku"`
	ProductName  string  `json:"product_name"`
	ProductTitle string  `json:"product_title"`
	VariantID    *string `json:"variant_id"`

	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	Currency   string  `json:"currency"`

	Discount float64  `json:"discount"`
	Tax      float64  `json:"tax"`
	TaxRate  *float64 `json:"tax_rate"`

	ImageURL          *string        `json:"image_url"`
	ProductURL        *string        `json:"product_url"`
	Weight            *float64       `json:"weight"`
	RequiresShipping  bool           `json:"requires_shipping"`
	IsGiftCard        bool           `json:"is_gift_card"`
	FulfillmentStatus *string        `json:"fulfillment_status"`
	Metadata          datatypes.JSON `json:"metadata"`
}

// Address representa una dirección en el dominio
type Address struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	Type    string `json:"type"`
	OrderID string `json:"order_id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
	Phone     string `json:"phone"`

	Street     string `json:"street"`
	Street2    string `json:"street2"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`

	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`

	Instructions *string        `json:"instructions"`
	IsDefault    bool           `json:"is_default"`
	Metadata     datatypes.JSON `json:"metadata"`
}

// Payment representa un pago en el dominio
type Payment struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID         string `json:"order_id"`
	PaymentMethodID uint   `json:"payment_method_id"`

	Amount       float64  `json:"amount"`
	Currency     string   `json:"currency"`
	ExchangeRate *float64 `json:"exchange_rate"`

	Status      string     `json:"status"`
	PaidAt      *time.Time `json:"paid_at"`
	ProcessedAt *time.Time `json:"processed_at"`

	TransactionID    *string `json:"transaction_id"`
	PaymentReference *string `json:"payment_reference"`
	Gateway          *string `json:"gateway"`

	RefundAmount  *float64       `json:"refund_amount"`
	RefundedAt    *time.Time     `json:"refunded_at"`
	FailureReason *string        `json:"failure_reason"`
	Metadata      datatypes.JSON `json:"metadata"`
}

// Shipment representa un envío en el dominio
type Shipment struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID string `json:"order_id"`

	TrackingNumber *string `json:"tracking_number"`
	TrackingURL    *string `json:"tracking_url"`
	Carrier        *string `json:"carrier"`
	CarrierCode    *string `json:"carrier_code"`

	GuideID  *string `json:"guide_id"`
	GuideURL *string `json:"guide_url"`

	Status      string     `json:"status"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`

	ShippingAddressID *uint `json:"shipping_address_id"`

	ShippingCost  *float64 `json:"shipping_cost"`
	InsuranceCost *float64 `json:"insurance_cost"`
	TotalCost     *float64 `json:"total_cost"`

	Weight *float64 `json:"weight"`
	Height *float64 `json:"height"`
	Width  *float64 `json:"width"`
	Length *float64 `json:"length"`

	WarehouseID   *uint  `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	DriverID      *uint  `json:"driver_id"`
	DriverName    string `json:"driver_name"`
	IsLastMile    bool   `json:"is_last_mile"`

	EstimatedDelivery *time.Time     `json:"estimated_delivery"`
	DeliveryNotes     *string        `json:"delivery_notes"`
	Metadata          datatypes.JSON `json:"metadata"`
}

// OrderChannelMetadata representa metadata del canal en el dominio
type OrderChannelMetadata struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	OrderID string `json:"order_id"`

	ChannelSource string `json:"channel_source"`
	IntegrationID uint   `json:"integration_id"`

	RawData datatypes.JSON `json:"raw_data"`

	Version     string     `json:"version"`
	ReceivedAt  time.Time  `json:"received_at"`
	ProcessedAt *time.Time `json:"processed_at"`
	IsLatest    bool       `json:"is_latest"`

	LastSyncedAt *time.Time `json:"last_synced_at"`
	SyncStatus   string     `json:"sync_status"`
}
