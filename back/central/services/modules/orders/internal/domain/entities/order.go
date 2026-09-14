package entities

import "time"

type ProbabilityOrder struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	BusinessID         *uint
	BusinessName       string
	IntegrationID      uint
	IntegrationType    string
	IntegrationLogoURL *string
	IntegrationName    string

	Platform       string
	ExternalID     string
	ChannelPackID  string
	OrderNumber    string
	InternalNumber string

	Subtotal              float64
	Tax                   float64
	Discount              float64
	ShippingCost          float64
	FreeShipping          bool
	TotalAmount           float64
	Currency              string
	IsCod                 bool
	CodTotal              *float64
	CodIncludesShipping   bool
	CodCheckoutCarrierFee float64

	SubtotalPresentment         float64
	TaxPresentment              float64
	DiscountPresentment         float64
	ShippingCostPresentment     float64
	ShippingDiscount            float64
	ShippingDiscountPresentment float64
	TotalAmountPresentment      float64
	CurrencyPresentment         string

	CustomerID        *uint
	CustomerName      string
	CustomerFirstName string
	CustomerLastName  string
	CustomerEmail     string
	CustomerPhone     string
	CustomerDNI       string

	ShippingStreet        string
	ShippingCity          string
	ShippingState         string
	ShippingCountry       string
	ShippingPostalCode    string
	ShippingLat           *float64
	ShippingLng           *float64
	ShippingGeoConfidence string

	ShippingAddressSource    string
	ShippingDeliveryType     string
	ShippingOfficeCarrier    string
	ShippingOfficeName       string
	ShippingOfficeDetails    []byte
	ShippingNeighborhood     string
	ShippingComplementType   string
	ShippingComplementNumber string
	ShippingTower            string
	ShippingBuilding         string
	DestinationDaneCode      string

	PaymentMethodID   uint
	PaymentMethodName string
	IsPaid            bool
	PaidAt            *time.Time

	TrackingNumber      *string
	TrackingLink        *string
	GuideID             *string
	GuideLink           *string
	DeliveryDate        *time.Time
	DeliveredAt         *time.Time
	DeliveryProbability *float64

	WarehouseID   *uint
	WarehouseName string
	DriverID      *uint
	DriverName    string
	IsLastMile    bool

	Weight *float64
	Height *float64
	Width  *float64
	Length *float64
	Boxes  *string

	OrderTypeID     *uint
	OrderTypeName   string
	Status          string
	StatusSource    string
	StatusChangedBy string
	StatusChangedAt *time.Time
	OriginalStatus  string
	StatusID        *uint
	OrderStatus     *OrderStatusInfo

	PaymentStatusID     *uint
	FulfillmentStatusID *uint
	PaymentStatus       *PaymentStatusInfo
	FulfillmentStatus   *FulfillmentStatusInfo

	Notes    *string
	Coupon   *string
	Approved *bool
	UserID   *uint
	UserName string

	UpdatedBy     *uint
	UpdatedByName string

	IsConfirmed *bool
	Novelty     *string

	IsTest bool

	Invoiceable     bool
	InvoiceURL      *string
	InvoiceID       *string
	InvoiceProvider *string
	InvoiceStatus   string // "", "pending", "issued", "failed", "cancelled"

	CodCutConfirmed bool

	OrderStatusURL string

	Metadata           []byte
	FinancialDetails   []byte
	ShippingDetails    []byte
	PaymentDetails     []byte
	FulfillmentDetails []byte

	OccurredAt time.Time
	ImportedAt time.Time

	OrderItems      []ProbabilityOrderItem
	Addresses       []ProbabilityAddress
	Payments        []ProbabilityPayment
	Shipments       []ProbabilityShipment
	ChannelMetadata []ProbabilityOrderChannelMetadata
	NegativeFactors []byte
	ScoreBreakdown  []byte

	CustomerOrderCount int
	CustomerTotalSpent string
	Address2           string
}

type OrderStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
}

type PaymentStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
}

type FulfillmentStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
}
