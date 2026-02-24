package request

import (
	"encoding/json"
)

// SerializableProbabilityOrderDTO es la estructura con etiquetas JSON para serialización.
// Esta estructura está en la capa de infraestructura, no en el dominio.
type SerializableProbabilityOrderDTO struct {
	BusinessID         *uint                           `json:"business_id"`
	IntegrationID      uint                            `json:"integration_id"`
	IntegrationType    string                          `json:"integration_type"`
	Platform           string                          `json:"platform"`
	ExternalID         string                          `json:"external_id"`
	OrderNumber        string                          `json:"order_number"`
	InternalNumber     string                          `json:"internal_number"`
	Subtotal           float64                         `json:"subtotal"`
	Tax                float64                         `json:"tax"`
	Discount           float64                         `json:"discount"`
	ShippingCost       float64                         `json:"shipping_cost"`
	TotalAmount        float64                         `json:"total_amount"`
	Currency           string                          `json:"currency"`
	CodTotal           *float64                        `json:"cod_total"`
	CustomerID         *uint                           `json:"customer_id"`
	CustomerName       string                          `json:"customer_name"`
	CustomerEmail      string                          `json:"customer_email"`
	CustomerPhone      string                          `json:"customer_phone"`
	CustomerDNI        string                          `json:"customer_dni"`
	OrderTypeID        *uint                           `json:"order_type_id"`
	OrderTypeName      string                          `json:"order_type_name"`
	Status             string                          `json:"status"`
	OriginalStatus     string                          `json:"original_status"`
	Notes              *string                         `json:"notes"`
	Coupon             *string                         `json:"coupon"`
	Approved           *bool                           `json:"approved"`
	UserID             *uint                           `json:"user_id"`
	UserName           string                          `json:"user_name"`
	Invoiceable        bool                            `json:"invoiceable"`
	InvoiceURL         *string                         `json:"invoice_url"`
	InvoiceID          *string                         `json:"invoice_id"`
	InvoiceProvider    *string                         `json:"invoice_provider"`
	OrderStatusURL     string                          `json:"order_status_url,omitempty"`
	OccurredAt         string                          `json:"occurred_at"`
	ImportedAt         string                          `json:"imported_at"`
	Items              json.RawMessage                 `json:"items,omitempty"`
	Metadata           json.RawMessage                 `json:"metadata,omitempty"`
	FinancialDetails   json.RawMessage                 `json:"financial_details,omitempty"`
	ShippingDetails    json.RawMessage                 `json:"shipping_details,omitempty"`
	PaymentDetails     json.RawMessage                 `json:"payment_details,omitempty"`
	FulfillmentDetails json.RawMessage                 `json:"fulfillment_details,omitempty"`
	OrderItems         []SerializableOrderItemDTO      `json:"order_items"`
	Addresses          []SerializableAddressDTO        `json:"addresses"`
	Payments           []SerializablePaymentDTO        `json:"payments"`
	Shipments          []SerializableShipmentDTO       `json:"shipments"`
	ChannelMetadata    *SerializableChannelMetadataDTO `json:"channel_metadata"`
}

type SerializableOrderItemDTO struct {
	ProductID    *string         `json:"product_id"`
	ProductSKU   string          `json:"product_sku"`
	ProductName  string          `json:"product_name"`
	ProductTitle string          `json:"product_title"`
	VariantID    *string         `json:"variant_id"`
	Quantity     int             `json:"quantity"`
	UnitPrice    float64         `json:"unit_price"`
	TotalPrice   float64         `json:"total_price"`
	Currency     string          `json:"currency"`
	Discount     float64         `json:"discount"`
	Tax          float64         `json:"tax"`
	TaxRate      *float64        `json:"tax_rate"`
	ImageURL     *string         `json:"image_url"`
	ProductURL   *string         `json:"product_url"`
	Weight       *float64        `json:"weight"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

type SerializableAddressDTO struct {
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
	Latitude     *float64        `json:"latitude"`
	Longitude    *float64        `json:"longitude"`
	Instructions *string         `json:"instructions"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

type SerializablePaymentDTO struct {
	PaymentMethodID  uint            `json:"payment_method_id"`
	Amount           float64         `json:"amount"`
	Currency         string          `json:"currency"`
	ExchangeRate     *float64        `json:"exchange_rate"`
	Status           string          `json:"status"`
	PaidAt           *string         `json:"paid_at"`
	ProcessedAt      *string         `json:"processed_at"`
	TransactionID    *string         `json:"transaction_id"`
	PaymentReference *string         `json:"payment_reference"`
	Gateway          *string         `json:"gateway"`
	RefundAmount     *float64        `json:"refund_amount"`
	RefundedAt       *string         `json:"refunded_at"`
	FailureReason    *string         `json:"failure_reason"`
	Metadata         json.RawMessage `json:"metadata,omitempty"`
}

type SerializableShipmentDTO struct {
	TrackingNumber    *string         `json:"tracking_number"`
	TrackingURL       *string         `json:"tracking_url"`
	Carrier           *string         `json:"carrier"`
	CarrierCode       *string         `json:"carrier_code"`
	GuideID           *string         `json:"guide_id"`
	GuideURL          *string         `json:"guide_url"`
	Status            string          `json:"status"`
	ShippedAt         *string         `json:"shipped_at"`
	DeliveredAt       *string         `json:"delivered_at"`
	ShippingAddressID *uint           `json:"shipping_address_id"`
	ShippingCost      *float64        `json:"shipping_cost"`
	InsuranceCost     *float64        `json:"insurance_cost"`
	TotalCost         *float64        `json:"total_cost"`
	Weight            *float64        `json:"weight"`
	Height            *float64        `json:"height"`
	Width             *float64        `json:"width"`
	Length            *float64        `json:"length"`
	WarehouseID       *uint           `json:"warehouse_id"`
	WarehouseName     string          `json:"warehouse_name"`
	DriverID          *uint           `json:"driver_id"`
	DriverName        string          `json:"driver_name"`
	IsLastMile        bool            `json:"is_last_mile"`
	EstimatedDelivery *string         `json:"estimated_delivery"`
	DeliveryNotes     *string         `json:"delivery_notes"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
}

type SerializableChannelMetadataDTO struct {
	ChannelSource string          `json:"channel_source"`
	RawData       json.RawMessage `json:"raw_data"`
	Version       string          `json:"version"`
	ReceivedAt    string          `json:"received_at"`
	ProcessedAt   *string         `json:"processed_at"`
	IsLatest      bool            `json:"is_latest"`
	LastSyncedAt  *string         `json:"last_synced_at"`
	SyncStatus    string          `json:"sync_status"`
}
