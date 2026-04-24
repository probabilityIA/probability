package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type InventoryLot struct {
	gorm.Model
	BusinessID      uint   `gorm:"not null;index;uniqueIndex:idx_lot_business_product_code,priority:1"`
	ProductID       string `gorm:"type:varchar(64);not null;index;uniqueIndex:idx_lot_business_product_code,priority:2"`
	LotCode         string `gorm:"size:100;not null;uniqueIndex:idx_lot_business_product_code,priority:3"`
	ManufactureDate *time.Time
	ExpirationDate  *time.Time     `gorm:"index"`
	ReceivedAt      *time.Time     `gorm:"index"`
	SupplierID      *uint          `gorm:"index"`
	Status          string         `gorm:"size:20;default:'active';index"`
	Metadata        datatypes.JSON `gorm:"type:jsonb"`

	Product  Product  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (InventoryLot) TableName() string {
	return "inventory_lots"
}

type InventorySerial struct {
	gorm.Model
	BusinessID         uint   `gorm:"not null;index;uniqueIndex:idx_serial_business_product_serial,priority:1"`
	ProductID          string `gorm:"type:varchar(64);not null;index;uniqueIndex:idx_serial_business_product_serial,priority:2"`
	SerialNumber       string `gorm:"size:100;not null;uniqueIndex:idx_serial_business_product_serial,priority:3"`
	LotID              *uint  `gorm:"index"`
	CurrentLocationID  *uint  `gorm:"index"`
	CurrentStateID     *uint  `gorm:"index"`
	ReceivedAt         *time.Time
	SoldAt             *time.Time

	Product  Product        `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business       `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Lot      *InventoryLot  `gorm:"foreignKey:LotID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (InventorySerial) TableName() string {
	return "inventory_serials"
}
