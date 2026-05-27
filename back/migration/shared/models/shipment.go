package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Shipment struct {
	gorm.Model

	OrderID *string `gorm:"type:varchar(36);index"`

	ClientName         string `gorm:"size:255"`
	DestinationAddress string `gorm:"size:255"`
	DestinationCity    string `gorm:"size:128"`
	DestinationState   string `gorm:"size:128"`
	DestinationSuburb  string `gorm:"size:128"`

	TrackingNumber *string `gorm:"size:128;index"`
	TrackingURL    *string `gorm:"size:512"`
	Carrier        *string `gorm:"size:128"`
	CarrierCode    *string `gorm:"size:50"`

	GuideID             *string `gorm:"size:128;index"`
	GuideURL            *string `gorm:"size:512"`
	ProbabilityGuideURL *string `gorm:"size:512"`

	Status              string     `gorm:"size:64;not null;index;default:'pending'"`
	CarrierStatus       *string    `gorm:"size:128;index"`
	CarrierStatusDetail *string    `gorm:"size:255"`
	ShippedAt           *time.Time `gorm:"index"`
	DeliveredAt         *time.Time

	ShippingAddressID *uint
	ShippingAddress   *Address `gorm:"foreignKey:ShippingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	DestinationGeozoneID   *uint          `gorm:"index"`
	DestinationGeozonePath datatypes.JSON `gorm:"type:jsonb"`
	GeozoneCountryID       *uint          `gorm:"index"`
	GeozoneStateID         *uint          `gorm:"index"`
	GeozoneCityID          *uint          `gorm:"index"`
	GeozoneAdminDistrictID *uint          `gorm:"index"`
	GeozoneLocalityID      *uint          `gorm:"index"`
	GeozoneNeighborhoodID  *uint          `gorm:"index"`
	GeozoneBarrioID        *uint          `gorm:"index"`

	ShippingCost         *float64 `gorm:"type:decimal(12,2)"`
	InsuranceCost        *float64 `gorm:"type:decimal(12,2)"`
	TotalCost            *float64 `gorm:"type:decimal(12,2)"`
	CarrierCost          *float64 `gorm:"type:decimal(12,2)"`
	AppliedMargin        *float64 `gorm:"type:decimal(12,2)"`
	CodCarrierFee        *float64 `gorm:"type:decimal(12,2)"`
	CodProbabilityMargin *float64 `gorm:"type:decimal(12,2)"`

	Weight *float64 `gorm:"type:decimal(10,2)"`
	Height *float64 `gorm:"type:decimal(10,2)"`
	Width  *float64 `gorm:"type:decimal(10,2)"`
	Length *float64 `gorm:"type:decimal(10,2)"`

	WarehouseID   *uint  `gorm:"index"`
	WarehouseName string `gorm:"size:128"`
	DriverID      *uint  `gorm:"index"`
	DriverName    string `gorm:"size:255"`
	IsLastMile    bool   `gorm:"default:false"`
	IsTest        bool   `gorm:"default:false;index"`

	EstimatedDelivery *time.Time     `gorm:"index"`
	DeliveryNotes     *string        `gorm:"type:text"`
	Metadata          datatypes.JSON `gorm:"type:jsonb"`

	Order *Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Shipment) TableName() string {
	return "shipments"
}
