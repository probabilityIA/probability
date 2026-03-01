package dtos

import (
	"encoding/json"
	"time"
)

// ProbabilityOrderDTO representa la estructura de orden en la lógica de negocio que todas las integraciones
// deben enviar después de mapear sus datos específicos
// Nota: Requiere tags JSON para deserialización en consumer
type ProbabilityOrderDTO struct {
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

	// Precios en moneda presentment (presentment_money - moneda local)
	SubtotalPresentment     float64 `json:"subtotal_presentment"`
	TaxPresentment          float64 `json:"tax_presentment"`
	DiscountPresentment     float64 `json:"discount_presentment"`
	ShippingCostPresentment float64 `json:"shipping_cost_presentment"`
	TotalAmountPresentment  float64 `json:"total_amount_presentment"`
	CurrencyPresentment     string  `json:"currency_presentment"`

	// Información del cliente
	CustomerID         *uint   `json:"customer_id,omitempty"`
	CustomerName       string  `json:"customer_name"`
	CustomerFirstName  string  `json:"customer_first_name"`
	CustomerLastName   string  `json:"customer_last_name"`
	CustomerEmail      string  `json:"customer_email"`
	CustomerPhone      string  `json:"customer_phone"`
	CustomerDNI        string  `json:"customer_dni"`
	CustomerOrderCount *int    `json:"customer_order_count,omitempty"`
	CustomerTotalSpent *string `json:"customer_total_spent,omitempty"`

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

	// Enlaces Externos
	OrderStatusURL string `json:"order_status_url"`

	// Timestamps
	OccurredAt time.Time `json:"occurred_at"`
	ImportedAt time.Time `json:"imported_at"`

	// Datos estructurados (JSONB) - Para compatibilidad
	// json.RawMessage permite deserializar JSON sin conocer la estructura exacta
	Items              json.RawMessage `json:"items,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	FinancialDetails   json.RawMessage `json:"financial_details,omitempty"`
	ShippingDetails    json.RawMessage `json:"shipping_details,omitempty"`
	PaymentDetails     json.RawMessage `json:"payment_details,omitempty"`
	FulfillmentDetails json.RawMessage `json:"fulfillment_details,omitempty"`

	// Control de flujo (no se persiste)
	IsManualOrder bool `json:"-"`

	// Tablas relacionadas
	OrderItems      []ProbabilityOrderItemDTO      `json:"order_items,omitempty"`
	Addresses       []ProbabilityAddressDTO        `json:"addresses,omitempty"`
	Payments        []ProbabilityPaymentDTO        `json:"payments,omitempty"`
	Shipments       []ProbabilityShipmentDTO       `json:"shipments,omitempty"`
	ChannelMetadata *ProbabilityChannelMetadataDTO `json:"channel_metadata,omitempty"`
}

// ProbabilityOrderItemDTO representa un item/producto de la orden
type ProbabilityOrderItemDTO struct {
	ProductID    *string  `json:"product_id,omitempty"`
	ProductSKU   string   `json:"product_sku"`
	ProductName  string   `json:"product_name"`
	ProductTitle string   `json:"product_title"`
	VariantID    *string  `json:"variant_id,omitempty"`
	Quantity     int      `json:"quantity"`
	UnitPrice    float64  `json:"unit_price"`
	TotalPrice   float64  `json:"total_price"`
	Currency     string   `json:"currency"`
	Discount     float64  `json:"discount"`
	Tax          float64  `json:"tax"`
	TaxRate      *float64 `json:"tax_rate,omitempty"`

	// Precios en moneda presentment (presentment_money - moneda local)
	UnitPricePresentment  float64         `json:"unit_price_presentment"`
	TotalPricePresentment float64         `json:"total_price_presentment"`
	DiscountPresentment   float64         `json:"discount_presentment"`
	TaxPresentment        float64         `json:"tax_presentment"`
	ImageURL              *string         `json:"image_url,omitempty"`
	ProductURL            *string         `json:"product_url,omitempty"`
	Weight                *float64        `json:"weight,omitempty"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
}

// ProbabilityAddressDTO representa una dirección (envío o facturación)
type ProbabilityAddressDTO struct {
	Type         string          `json:"type"`
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	Company      string          `json:"company"`
	Phone        string          `json:"phone"`
	Street       string          `json:"street"`
	Street2      string          `json:"street2"`
	City         string          `json:"city"`
	State        string          `json:"state"`
	Country      string          `json:"country"`
	PostalCode   string          `json:"postal_code"`
	Latitude     *float64        `json:"latitude,omitempty"`
	Longitude    *float64        `json:"longitude,omitempty"`
	Instructions *string         `json:"instructions,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// ProbabilityPaymentDTO representa un pago de la orden
type ProbabilityPaymentDTO struct {
	PaymentMethodID  uint            `json:"payment_method_id"`
	Amount           float64         `json:"amount"`
	Currency         string          `json:"currency"`
	ExchangeRate     *float64        `json:"exchange_rate,omitempty"`
	Status           string          `json:"status"`
	PaidAt           *time.Time      `json:"paid_at,omitempty"`
	ProcessedAt      *time.Time      `json:"processed_at,omitempty"`
	TransactionID    *string         `json:"transaction_id,omitempty"`
	PaymentReference *string         `json:"payment_reference,omitempty"`
	Gateway          *string         `json:"gateway,omitempty"`
	RefundAmount     *float64        `json:"refund_amount,omitempty"`
	RefundedAt       *time.Time      `json:"refunded_at,omitempty"`
	FailureReason    *string         `json:"failure_reason,omitempty"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
}

// ProbabilityShipmentDTO representa un envío de la orden
type ProbabilityShipmentDTO struct {
	TrackingNumber    *string         `json:"tracking_number,omitempty"`
	TrackingURL       *string         `json:"tracking_url,omitempty"`
	Carrier           *string         `json:"carrier,omitempty"`
	CarrierCode       *string         `json:"carrier_code,omitempty"`
	GuideID           *string         `json:"guide_id,omitempty"`
	GuideURL          *string         `json:"guide_url,omitempty"`
	Status            string          `json:"status"`
	ShippedAt         *time.Time      `json:"shipped_at,omitempty"`
	DeliveredAt       *time.Time      `json:"delivered_at,omitempty"`
	ShippingAddressID *uint           `json:"shipping_address_id,omitempty"`
	ShippingCost      *float64        `json:"shipping_cost,omitempty"`
	InsuranceCost     *float64        `json:"insurance_cost,omitempty"`
	TotalCost         *float64        `json:"total_cost,omitempty"`
	Weight            *float64        `json:"weight,omitempty"`
	Height            *float64        `json:"height,omitempty"`
	Width             *float64        `json:"width,omitempty"`
	Length            *float64        `json:"length,omitempty"`
	WarehouseID       *uint           `json:"warehouse_id,omitempty"`
	WarehouseName     string          `json:"warehouse_name"`
	DriverID          *uint           `json:"driver_id,omitempty"`
	DriverName        string          `json:"driver_name"`
	IsLastMile        bool            `json:"is_last_mile"`
	EstimatedDelivery *time.Time      `json:"estimated_delivery,omitempty"`
	DeliveryNotes     *string         `json:"delivery_notes,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
}

// ProbabilityChannelMetadataDTO representa los datos crudos del canal
type ProbabilityChannelMetadataDTO struct {
	ChannelSource string          `json:"channel_source"`
	RawData       json.RawMessage `json:"raw_data,omitempty"`
	Version       string          `json:"version"`
	ReceivedAt    time.Time       `json:"received_at"`
	ProcessedAt   *time.Time      `json:"processed_at,omitempty"`
	IsLatest      bool            `json:"is_latest"`
	LastSyncedAt  *time.Time      `json:"last_synced_at,omitempty"`
	SyncStatus    string          `json:"sync_status"`
}
