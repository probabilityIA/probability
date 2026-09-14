package canonical

import (
	"time"
)

type ProbabilityOrderDTO struct {
	BusinessID      *uint
	IntegrationID   uint
	IntegrationType string
	Platform        string
	ExternalID      string
	ChannelPackID   string
	OrderNumber     string
	InternalNumber  string
	Subtotal        float64
	Tax             float64
	Discount        float64
	ShippingCost    float64
	FreeShipping    bool
	TotalAmount     float64
	Currency        string
	CodTotal        *float64

	CodCheckoutCarrierFee float64

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
	Items              []byte
	Metadata           []byte
	FinancialDetails   []byte
	ShippingDetails    []byte
	PaymentDetails     []byte
	FulfillmentDetails []byte
	OrderItems         []ProbabilityOrderItemDTO
	Addresses          []ProbabilityAddressDTO
	Payments           []ProbabilityPaymentDTO
	Shipments          []ProbabilityShipmentDTO
	ChannelMetadata    *ProbabilityChannelMetadataDTO
	SkipInventory      bool
}

type ProbabilityOrderItemDTO struct {
	ProductID       *string
	ProductSKU      string
	ProductName     string
	ProductTitle    string
	VariantID       *string
	ExternalBarcode *string
	Quantity        int
	UnitPrice       float64
	TotalPrice      float64
	Currency        string
	Discount        float64
	Tax             float64
	TaxRate         *float64
	ImageURL        *string
	ProductURL      *string
	Weight          *float64
	Metadata        []byte

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
	Metadata         []byte
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
	Metadata          []byte
}

type ProbabilityChannelMetadataDTO struct {
	ChannelSource string
	RawData       []byte
	Version       string
	ReceivedAt    time.Time
	ProcessedAt   *time.Time
	IsLatest      bool
	LastSyncedAt  *time.Time
	SyncStatus    string
}
