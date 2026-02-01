package request

import (
	"time"

	"gorm.io/datatypes"
)

// MapOrder representa la petición HTTP para mapear una orden canónica
// ✅ DTO HTTP - CON TAGS (json + binding + datatypes.JSON)
type MapOrder struct {
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

	// Precios en moneda presentment (presentment_money - moneda local)
	SubtotalPresentment     float64 `json:"subtotal_presentment,omitempty"`
	TaxPresentment          float64 `json:"tax_presentment,omitempty"`
	DiscountPresentment     float64 `json:"discount_presentment,omitempty"`
	ShippingCostPresentment float64 `json:"shipping_cost_presentment,omitempty"`
	TotalAmountPresentment  float64 `json:"total_amount_presentment,omitempty"`
	CurrencyPresentment     string  `json:"currency_presentment,omitempty"`

	// Información del cliente
	CustomerID         *uint   `json:"customer_id"`
	CustomerName       string  `json:"customer_name" binding:"max=255"`
	CustomerEmail      string  `json:"customer_email" binding:"max=255"`
	CustomerPhone      string  `json:"customer_phone" binding:"max=32"`
	CustomerDNI        string  `json:"customer_dni" binding:"max=64"`
	CustomerOrderCount *int    `json:"customer_order_count"`
	CustomerTotalSpent *string `json:"customer_total_spent"`

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

	// Enlaces Externos
	OrderStatusURL string `json:"order_status_url,omitempty"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`

	// Datos estructurados (JSONB) - permitido usar datatypes.JSON en infra
	Items              datatypes.JSON `json:"items,omitempty"`
	Metadata           datatypes.JSON `json:"metadata,omitempty"`
	FinancialDetails   datatypes.JSON `json:"financial_details,omitempty"`
	ShippingDetails    datatypes.JSON `json:"shipping_details,omitempty"`
	PaymentDetails     datatypes.JSON `json:"payment_details,omitempty"`
	FulfillmentDetails datatypes.JSON `json:"fulfillment_details,omitempty"`

	// Tablas relacionadas
	OrderItems      []MapOrderItem      `json:"order_items" binding:"dive"`
	Addresses       []MapAddress        `json:"addresses" binding:"dive"`
	Payments        []MapPayment        `json:"payments" binding:"dive"`
	Shipments       []MapShipment       `json:"shipments" binding:"dive"`
	ChannelMetadata *MapChannelMetadata `json:"channel_metadata"`
}

// MapOrderItem representa un item de la orden en HTTP
type MapOrderItem struct {
	ProductID    *string  `json:"product_id"`
	ProductSKU   string   `json:"product_sku" binding:"required,max=128"`
	ProductName  string   `json:"product_name" binding:"required,max=255"`
	ProductTitle string   `json:"product_title" binding:"max=255"`
	VariantID    *string  `json:"variant_id"`
	Quantity     int      `json:"quantity" binding:"required,min=1"`
	UnitPrice    float64  `json:"unit_price" binding:"required,min=0"`
	TotalPrice   float64  `json:"total_price" binding:"required,min=0"`
	Currency     string   `json:"currency" binding:"max=10"`
	Discount     float64  `json:"discount" binding:"min=0"`
	Tax          float64  `json:"tax" binding:"min=0"`
	TaxRate      *float64 `json:"tax_rate"`

	// Precios en moneda presentment
	UnitPricePresentment  float64        `json:"unit_price_presentment,omitempty"`
	TotalPricePresentment float64        `json:"total_price_presentment,omitempty"`
	DiscountPresentment   float64        `json:"discount_presentment,omitempty"`
	TaxPresentment        float64        `json:"tax_presentment,omitempty"`
	ImageURL              *string        `json:"image_url"`
	ProductURL            *string        `json:"product_url"`
	Weight                *float64       `json:"weight"`
	Metadata              datatypes.JSON `json:"metadata,omitempty"`
}

// MapAddress representa una dirección en HTTP
type MapAddress struct {
	Type         string         `json:"type" binding:"required,oneof=shipping billing"`
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

// MapPayment representa un pago en HTTP
type MapPayment struct {
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

// MapShipment representa un envío en HTTP
type MapShipment struct {
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

// MapChannelMetadata representa metadata del canal en HTTP
type MapChannelMetadata struct {
	ChannelSource string         `json:"channel_source" binding:"required,max=50"`
	RawData       datatypes.JSON `json:"raw_data" binding:"required"`
	Version       string         `json:"version" binding:"max=20"`
	ReceivedAt    time.Time      `json:"received_at"`
	ProcessedAt   *time.Time     `json:"processed_at"`
	IsLatest      bool           `json:"is_latest"`
	LastSyncedAt  *time.Time     `json:"last_synced_at"`
	SyncStatus    string         `json:"sync_status" binding:"max=64"`
}
