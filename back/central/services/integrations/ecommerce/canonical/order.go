// Package canonical define el formato canónico de órdenes para todos los módulos
// de integración e-commerce. Todos los providers (Shopify, MercadoLibre, WooCommerce, etc.)
// deben mapear sus órdenes a este formato antes de publicarlas a la cola de RabbitMQ.
//
// No tiene etiquetas JSON — es una estructura de dominio pura.
// La serialización (con etiquetas JSON) se realiza en la capa de infraestructura
// de cada módulo (infra/secondary/queue/request/).
package canonical

import (
	"time"

	"gorm.io/datatypes"
)

// ProbabilityOrderDTO es el formato canónico de orden que todos los módulos
// de e-commerce publican a la cola probability.orders.canonical.
type ProbabilityOrderDTO struct {
	BusinessID      *uint
	IntegrationID   uint
	IntegrationType string
	Platform        string
	ExternalID      string
	OrderNumber     string
	InternalNumber  string
	Subtotal        float64
	Tax             float64
	Discount        float64
	ShippingCost    float64
	TotalAmount     float64
	Currency        string
	CodTotal        *float64

	// Precios en moneda presentment (moneda local del cliente)
	SubtotalPresentment     float64
	TaxPresentment          float64
	DiscountPresentment     float64
	ShippingCostPresentment float64
	TotalAmountPresentment  float64
	CurrencyPresentment     string

	CustomerID         *uint
	CustomerName       string
	CustomerEmail      string
	CustomerPhone      string
	CustomerDNI        string
	CustomerOrderCount *int
	CustomerTotalSpent *string
	OrderTypeID        *uint
	OrderTypeName      string
	Status             string
	OriginalStatus     string
	Notes              *string
	Coupon             *string
	Approved           *bool
	UserID             *uint
	UserName           string
	Invoiceable        bool
	InvoiceURL         *string
	InvoiceID          *string
	InvoiceProvider    *string
	OrderStatusURL     string
	OccurredAt         time.Time
	ImportedAt         time.Time
	Items              datatypes.JSON
	Metadata           datatypes.JSON
	FinancialDetails   datatypes.JSON
	ShippingDetails    datatypes.JSON
	PaymentDetails     datatypes.JSON
	FulfillmentDetails datatypes.JSON
	OrderItems         []ProbabilityOrderItemDTO
	Addresses          []ProbabilityAddressDTO
	Payments           []ProbabilityPaymentDTO
	Shipments          []ProbabilityShipmentDTO
	ChannelMetadata    *ProbabilityChannelMetadataDTO
}

type ProbabilityOrderItemDTO struct {
	ProductID    *string
	ProductSKU   string
	ProductName  string
	ProductTitle string
	VariantID    *string
	Quantity     int
	UnitPrice    float64
	TotalPrice   float64
	Currency     string
	Discount     float64
	Tax          float64
	TaxRate      *float64
	ImageURL     *string
	ProductURL   *string
	Weight       *float64
	Metadata     []byte

	// Precios en moneda presentment
	UnitPricePresentment  float64
	TotalPricePresentment float64
	DiscountPresentment   float64
	TaxPresentment        float64
}

type ProbabilityAddressDTO struct {
	Type         string
	FirstName    string
	LastName     string
	Company      string
	Phone        string
	Street       string
	Street2      string
	City         string
	State        string
	Country      string
	PostalCode   string
	Latitude     *float64
	Longitude    *float64
	Instructions *string
	Metadata     []byte
}

type ProbabilityPaymentDTO struct {
	PaymentMethodID  uint
	Amount           float64
	Currency         string
	ExchangeRate     *float64
	Status           string
	PaidAt           *time.Time
	ProcessedAt      *time.Time
	TransactionID    *string
	PaymentReference *string
	Gateway          *string
	RefundAmount     *float64
	RefundedAt       *time.Time
	FailureReason    *string
	Metadata         datatypes.JSON
}

type ProbabilityShipmentDTO struct {
	TrackingNumber    *string
	TrackingURL       *string
	Carrier           *string
	CarrierCode       *string
	GuideID           *string
	GuideURL          *string
	Status            string
	ShippedAt         *time.Time
	DeliveredAt       *time.Time
	ShippingAddressID *uint
	ShippingCost      *float64
	InsuranceCost     *float64
	TotalCost         *float64
	Weight            *float64
	Height            *float64
	Width             *float64
	Length            *float64
	WarehouseID       *uint
	WarehouseName     string
	DriverID          *uint
	DriverName        string
	IsLastMile        bool
	EstimatedDelivery *time.Time
	DeliveryNotes     *string
	Metadata          datatypes.JSON
}

type ProbabilityChannelMetadataDTO struct {
	ChannelSource string
	RawData       datatypes.JSON
	Version       string
	ReceivedAt    time.Time
	ProcessedAt   *time.Time
	IsLatest      bool
	LastSyncedAt  *time.Time
	SyncStatus    string
}
