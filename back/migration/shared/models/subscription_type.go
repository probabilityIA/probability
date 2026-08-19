package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SubscriptionType struct {
	gorm.Model
	Name                 string         `gorm:"size:100;not null"`
	Code                 string         `gorm:"size:50;not null;unique"`
	Description          string         `gorm:"type:text"`
	Price                float64        `gorm:"type:numeric;not null"`
	BillingPeriod        string         `gorm:"size:20;not null;default:'monthly'"`
	Active               bool           `gorm:"default:true"`
	Features             datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'"`
	MaxEcommerceChannels int            `gorm:"not null;default:0"`
	BusinessID           *uint          `gorm:"index"`
	IncludedShipments    *int           `gorm:"column:included_shipments"`
	ShipmentOveragePrice *float64       `gorm:"column:shipment_overage_price;type:numeric"`
	IncludedInvoices     *int           `gorm:"column:included_invoices"`
	InvoiceOveragePrice  *float64       `gorm:"column:invoice_overage_price;type:numeric"`
	IncludedOrders       *int           `gorm:"column:included_orders"`
	OrderOveragePrice    *float64       `gorm:"column:order_overage_price;type:numeric"`
}

func (SubscriptionType) TableName() string {
	return "subscription_types"
}
