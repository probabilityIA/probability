package models

import (
	"crypto/rand"
	"fmt"
	mathrand "math/rand"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Order struct {
	ID        string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	BusinessID      *uint  `gorm:"index"`
	IntegrationID   uint   `gorm:"not null;index;uniqueIndex:idx_integration_external_id,priority:1"`
	IntegrationType string `gorm:"size:50;not null;index"`

	Platform       string `gorm:"size:50;not null;index"`
	ExternalID     string `gorm:"size:255;not null;index;uniqueIndex:idx_integration_external_id,priority:2"`
	OrderNumber    string `gorm:"size:128;index"`
	InternalNumber string `gorm:"size:128;unique;index"`

	Subtotal     float64  `gorm:"type:decimal(12,2);not null;default:0"`
	Tax          float64  `gorm:"type:decimal(12,2);not null;default:0"`
	Discount     float64  `gorm:"type:decimal(12,2);not null;default:0"`
	ShippingCost float64  `gorm:"type:decimal(12,2);not null;default:0"`
	TotalAmount  float64  `gorm:"type:decimal(12,2);not null"`
	Currency     string   `gorm:"size:10;default:'USD'"`
	CodTotal     *float64 `gorm:"type:decimal(12,2)"`

	SubtotalPresentment         float64 `gorm:"column:subtotal_presentment;type:decimal(12,2);not null;default:0"`
	TaxPresentment              float64 `gorm:"column:tax_presentment;type:decimal(12,2);not null;default:0"`
	DiscountPresentment         float64 `gorm:"column:discount_presentment;type:decimal(12,2);not null;default:0"`
	ShippingCostPresentment     float64 `gorm:"column:shipping_cost_presentment;type:decimal(12,2);not null;default:0"`
	ShippingDiscount            float64 `gorm:"column:shipping_discount;type:decimal(12,2);not null;default:0"`
	ShippingDiscountPresentment float64 `gorm:"column:shipping_discount_presentment;type:decimal(12,2);not null;default:0"`
	TotalAmountPresentment      float64 `gorm:"column:total_amount_presentment;type:decimal(12,2);not null;default:0"`
	CurrencyPresentment         string  `gorm:"size:10"`

	CustomerID        *uint  `gorm:"index"`
	CustomerName      string `gorm:"size:255"`
	CustomerFirstName string `gorm:"size:128"`
	CustomerLastName  string `gorm:"size:128"`
	CustomerEmail     string `gorm:"size:255;index"`
	CustomerPhone     string `gorm:"size:50"`
	CustomerDNI       string `gorm:"size:64"`

	ShippingStreet     string `gorm:"size:255"`
	ShippingCity       string `gorm:"size:128"`
	ShippingState      string `gorm:"size:128"`
	ShippingCountry    string `gorm:"size:128"`
	ShippingPostalCode string `gorm:"size:32"`
	ShippingLat        *float64
	ShippingLng        *float64

	ShippingGeoConfidence string `gorm:"size:16;index"`

	DestinationGeozoneID   *uint          `gorm:"index"`
	DestinationGeozonePath datatypes.JSON `gorm:"type:jsonb"`
	GeozoneCountryID       *uint          `gorm:"index"`
	GeozoneStateID         *uint          `gorm:"index"`
	GeozoneCityID          *uint          `gorm:"index"`
	GeozoneAdminDistrictID *uint          `gorm:"index"`
	GeozoneLocalityID      *uint          `gorm:"index"`
	GeozoneNeighborhoodID  *uint          `gorm:"index"`
	GeozoneBarrioID        *uint          `gorm:"index"`

	PaymentMethodID uint `gorm:"not null;index"`
	IsPaid          bool `gorm:"default:false;index"`
	PaidAt          *time.Time

	TrackingNumber      *string    `gorm:"size:128;index"`
	TrackingLink        *string    `gorm:"size:512"`
	GuideID             *string    `gorm:"size:128"`
	GuideLink           *string    `gorm:"size:512"`
	DeliveryDate        *time.Time `gorm:"index"`
	DeliveredAt         *time.Time
	DeliveryProbability *float64 `gorm:"type:decimal(5,2)"`

	WarehouseID   *uint  `gorm:"index"`
	WarehouseName string `gorm:"size:128"`
	DriverID      *uint  `gorm:"index"`
	DriverName    string `gorm:"size:255"`
	IsLastMile    bool   `gorm:"default:false"`

	Weight *float64 `gorm:"type:decimal(10,2)"`
	Height *float64 `gorm:"type:decimal(10,2)"`
	Width  *float64 `gorm:"type:decimal(10,2)"`
	Length *float64 `gorm:"type:decimal(10,2)"`
	Boxes  *string  `gorm:"type:text"`

	OrderTypeID    *uint  `gorm:"index"`
	OrderTypeName  string `gorm:"size:64"`
	Status         string `gorm:"size:64;not null;index;default:'pending'"`
	OriginalStatus string `gorm:"size:64"`
	StatusID       *uint  `gorm:"index"`

	PaymentStatusID     *uint `gorm:"index"`
	FulfillmentStatusID *uint `gorm:"index"`

	Notes    *string `gorm:"type:text"`
	Coupon   *string `gorm:"size:128"`
	Approved *bool
	UserID   *uint  `gorm:"index"`
	UserName string `gorm:"size:255"`

	IsConfirmed *bool   `gorm:"default:false"`
	Novelty     *string `gorm:"type:text"`

	IsTest bool `gorm:"default:false;index"`

	Invoiceable     bool    `gorm:"default:false"`
	InvoiceURL      *string `gorm:"size:512"`
	InvoiceID       *string `gorm:"size:128;index"`
	InvoiceProvider *string `gorm:"size:64"`

	OrderStatusURL string `gorm:"size:512"`

	Metadata datatypes.JSON `gorm:"type:jsonb"`

	NegativeFactors datatypes.JSON `gorm:"type:jsonb"`

	ScoreBreakdown datatypes.JSON `gorm:"type:jsonb"`

	FinancialDetails datatypes.JSON `gorm:"type:jsonb"`

	ShippingDetails datatypes.JSON `gorm:"type:jsonb"`

	PaymentDetails datatypes.JSON `gorm:"type:jsonb"`

	FulfillmentDetails datatypes.JSON `gorm:"type:jsonb"`

	OccurredAt time.Time `gorm:"index"`
	ImportedAt time.Time `gorm:"index"`

	Business          *Business         `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Integration       Integration       `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PaymentMethod     PaymentMethod     `gorm:"foreignKey:PaymentMethodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderStatus       OrderStatus       `gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PaymentStatus     PaymentStatus     `gorm:"foreignKey:PaymentStatusID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	FulfillmentStatus FulfillmentStatus `gorm:"foreignKey:FulfillmentStatusID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	OrderItems      []OrderItem            `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Addresses       []Address              `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Payments        []Payment              `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Shipments       []Shipment             `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ChannelMetadata []OrderChannelMetadata `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Order) TableName() string {
	return "orders"
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {

	if o.ID == "" {
		o.ID = generateUUID()
	}

	if o.InternalNumber == "" {
		o.InternalNumber = fmt.Sprintf("ORD-%d-%s",
			time.Now().Unix(),
			generateRandomString(6))
	}

	if o.CreatedAt.IsZero() && !o.OccurredAt.IsZero() {
		o.CreatedAt = o.OccurredAt
	}

	return nil
}

type OrderHistory struct {
	gorm.Model
	OrderID        string         `gorm:"type:varchar(36);not null;index"`
	PreviousStatus string         `gorm:"size:64"`
	NewStatus      string         `gorm:"size:64;not null"`
	ChangedBy      *uint          `gorm:"index"`
	ChangedByName  string         `gorm:"size:255"`
	Reason         *string        `gorm:"type:text"`
	Metadata       datatypes.JSON `gorm:"type:jsonb"`

	Order Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (OrderHistory) TableName() string {
	return "order_history"
}

func generateUUID() string {

	b := make([]byte, 16)
	rand.Read(b)

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func generateRandomString(n int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

type OrderItem struct {
	gorm.Model

	OrderID   string  `gorm:"type:varchar(36);not null;index"`
	ProductID *string `gorm:"type:varchar(64);index"`
	Order     Order   `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Quantity   int     `gorm:"not null;default:1"`
	UnitPrice  float64 `gorm:"type:decimal(12,2);not null"`
	TotalPrice float64 `gorm:"type:decimal(12,2);not null"`
	Currency   string  `gorm:"size:10;default:'USD'"`

	Discount        float64  `gorm:"type:decimal(12,2);default:0"`
	DiscountPercent float64  `gorm:"type:decimal(5,2);default:0"`
	Tax             float64  `gorm:"type:decimal(12,2);default:0"`
	TaxRate         *float64 `gorm:"type:decimal(5,4)"`

	UnitPriceBase            float64 `gorm:"column:unit_price_base;type:decimal(12,2);not null;default:0"`
	UnitPriceBasePresentment float64 `gorm:"column:unit_price_base_presentment;type:decimal(12,2);not null;default:0"`

	UnitPricePresentment  float64 `gorm:"column:unit_price_presentment;type:decimal(12,2);not null;default:0"`
	TotalPricePresentment float64 `gorm:"column:total_price_presentment;type:decimal(12,2);not null;default:0"`
	DiscountPresentment   float64 `gorm:"column:discount_presentment;type:decimal(12,2);default:0"`
	TaxPresentment        float64 `gorm:"column:tax_presentment;type:decimal(12,2);default:0"`

	ProductSKU        string         `gorm:"size:255"`
	ProductName       string         `gorm:"size:255"`
	VariantID         *string        `gorm:"size:255"`
	VariantLabel      string         `gorm:"size:255"`
	FulfillmentStatus *string        `gorm:"size:64"`
	Metadata          datatypes.JSON `gorm:"type:jsonb"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type Address struct {
	gorm.Model

	Type string `gorm:"size:20;not null;index"`

	OrderID string `gorm:"type:varchar(36);not null;index"`

	FirstName string `gorm:"size:128"`
	LastName  string `gorm:"size:128"`
	Company   string `gorm:"size:255"`
	Phone     string `gorm:"size:32"`

	Street     string `gorm:"size:255;not null"`
	Street2    string `gorm:"size:255"`
	City       string `gorm:"size:128;not null"`
	State      string `gorm:"size:128"`
	Country    string `gorm:"size:128;not null"`
	PostalCode string `gorm:"size:32"`

	Latitude  *float64 `gorm:"type:decimal(10,8)"`
	Longitude *float64 `gorm:"type:decimal(11,8)"`

	Instructions *string        `gorm:"type:text"`
	IsDefault    bool           `gorm:"default:false"`
	Metadata     datatypes.JSON `gorm:"type:jsonb"`

	Order Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Address) TableName() string {
	return "addresses"
}

type Payment struct {
	gorm.Model

	OrderID string `gorm:"type:varchar(36);not null;index"`

	PaymentMethodID uint `gorm:"not null;index"`

	Amount       float64  `gorm:"type:decimal(12,2);not null"`
	Currency     string   `gorm:"size:10;default:'USD'"`
	ExchangeRate *float64 `gorm:"type:decimal(10,4)"`

	Status      string     `gorm:"size:64;not null;index"`
	PaidAt      *time.Time `gorm:"index"`
	ProcessedAt *time.Time

	TransactionID    *string `gorm:"size:255;index"`
	PaymentReference *string `gorm:"size:255"`
	Gateway          *string `gorm:"size:64"`

	RefundAmount  *float64 `gorm:"type:decimal(12,2)"`
	RefundedAt    *time.Time
	FailureReason *string        `gorm:"type:text"`
	Metadata      datatypes.JSON `gorm:"type:jsonb"`

	Order         Order         `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Payment) TableName() string {
	return "payments"
}

type OrderChannelMetadata struct {
	gorm.Model

	OrderID string `gorm:"type:varchar(36);not null;index"`

	ChannelSource string `gorm:"size:50;not null;index"`
	IntegrationID uint   `gorm:"not null;index"`

	RawData datatypes.JSON `gorm:"type:jsonb;not null"`

	Version     string    `gorm:"size:20"`
	ReceivedAt  time.Time `gorm:"index"`
	ProcessedAt *time.Time
	IsLatest    bool `gorm:"default:true;index"`

	LastSyncedAt *time.Time
	SyncStatus   string `gorm:"size:64;default:'pending'"`

	Order       Order       `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Integration Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (OrderChannelMetadata) TableName() string {
	return "order_channel_metadata"
}
