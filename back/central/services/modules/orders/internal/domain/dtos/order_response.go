package dtos

import (
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

type OrderResponse struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	BusinessID         *uint
	IntegrationID      uint
	IntegrationType    string
	IntegrationLogoURL *string
	IntegrationName    string

	Platform       string
	ExternalID     string
	ChannelPackID  string
	OrderNumber    string
	InternalNumber string

	Subtotal                    float64
	Tax                         float64
	Discount                    float64
	ShippingCost                float64
	ShippingDiscount            float64
	ShippingDiscountPresentment float64
	TotalAmount                 float64
	Currency                    string
	IsCod                       bool
	CodTotal                    *float64
	CodIncludesShipping         bool
	CodCutConfirmed             bool

	SubtotalPresentment     float64
	TaxPresentment          float64
	DiscountPresentment     float64
	ShippingCostPresentment float64
	TotalAmountPresentment  float64
	CurrencyPresentment     string

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

	PaymentMethodID uint
	IsPaid          bool
	PaidAt          *time.Time

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

	OrderTypeID    *uint
	OrderTypeName  string
	Status         string
	OriginalStatus string
	StatusID       *uint
	OrderStatus    *entities.OrderStatusInfo

	PaymentStatusID     *uint
	FulfillmentStatusID *uint
	PaymentStatus       *entities.PaymentStatusInfo
	FulfillmentStatus   *entities.FulfillmentStatusInfo

	Notes    *string
	Coupon   *string
	Approved *bool
	UserID   *uint
	UserName string

	IsConfirmed *bool
	Novelty     *string

	IsTest bool

	Invoiceable     bool
	InvoiceURL      *string
	InvoiceID       *string
	InvoiceProvider *string

	OrderStatusURL string

	OrderItems []entities.ProbabilityOrderItem

	Shipment *ShipmentData

	Invoice *InvoiceData

	Metadata           []byte
	FinancialDetails   []byte
	ShippingDetails    []byte
	FreeShipping       bool
	PaymentDetails     []byte
	FulfillmentDetails []byte
	NegativeFactors    []string
	ScoreBreakdown     json.RawMessage

	OccurredAt time.Time
	ImportedAt time.Time
}

type ShipmentData struct {
	ID                  uint
	Carrier             *string
	TrackingNumber      *string
	GuideURL            *string
	Status              string
	CarrierStatus       *string
	CarrierStatusDetail *string
	TotalCost           *float64
	CodCarrierFee       *float64
}

type InvoiceData struct {
	ID              uint
	InvoiceNumber   string
	Status          string
	IssuedAt        *time.Time
	RetentionAmount float64
}
