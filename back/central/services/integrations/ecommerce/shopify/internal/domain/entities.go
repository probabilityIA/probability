package domain

import (
	"time"

	core "github.com/secamc93/probability/back/central/services/integrations/core"
	"gorm.io/datatypes"
)

// Integration es un alias directo de core.PublicIntegration
type Integration = core.PublicIntegration

// IntegrationTypeID es el ID del tipo de integración Shopify en la base de datos.
const IntegrationTypeID = 1

type ShopifyOrderDTO struct {
	ID                string
	OrderNumber       string
	TotalPrice        float64
	Currency          string
	PaymentType       string
	CustomerName      string
	CustomerEmail     string
	Phone             string
	Country           string
	Province          string
	City              string
	Address           string
	AddressComplement string
	FinancialStatus   string
	FulfillmentStatus string
	CreatedAt         time.Time
	RawData           map[string]interface{}
}

type ShopifyOrder struct {
	BusinessID      *uint
	IntegrationID   uint
	IntegrationType string
	Platform        string
	ExternalID      string
	OrderNumber     string
	TotalAmount     float64
	Subtotal        float64
	Tax             float64
	Discount        float64
	ShippingCost    float64
	Currency        string
	Customer        ShopifyCustomer
	ShippingAddress ShopifyAddress
	Status          string
	OriginalStatus  string
	Items           []ShopifyOrderItem
	Metadata        map[string]interface{}
	OccurredAt      time.Time
	ImportedAt      time.Time
	OrderStatusURL  string
	RawData         []byte

	// Precios en moneda presentment (presentment_money - moneda local)
	SubtotalPresentment     float64
	TaxPresentment          float64
	DiscountPresentment     float64
	ShippingCostPresentment float64
	TotalAmountPresentment  float64
	CurrencyPresentment     string
}

type ShopifyCustomer struct {
	Name           string
	Email          string
	Phone          string
	DefaultAddress *ShopifyAddress
	OrdersCount    int
	TotalSpent     string
}

type ShopifyAddress struct {
	Street      string
	Address2    string
	City        string
	State       string
	Country     string
	PostalCode  string
	Coordinates *struct {
		Lat float64
		Lng float64
	}
}

type ShopifyOrderItem struct {
	ExternalID   string
	Name         string
	SKU          string
	Quantity     int
	UnitPrice    float64
	ProductID    *int64   // ID del producto en Shopify
	VariantID    *int64   // ID de la variante en Shopify
	Title        string   // Título del producto
	VariantTitle *string  // Título de la variante
	Discount     float64  // Descuento aplicado
	Tax          float64  // Impuesto
	Weight       *float64 // Peso en gramos

	// Precios en moneda presentment (presentment_money - moneda local)
	UnitPricePresentment float64
	DiscountPresentment  float64
	TaxPresentment       float64
}

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

	// Precios en moneda presentment (presentment_money - moneda local)
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

	// Precios en moneda presentment (presentment_money - moneda local)
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
	RawData       datatypes.JSON // Cambiado de []byte a datatypes.JSON para consistencia
	Version       string
	ReceivedAt    time.Time
	ProcessedAt   *time.Time
	IsLatest      bool
	LastSyncedAt  *time.Time
	SyncStatus    string
}
