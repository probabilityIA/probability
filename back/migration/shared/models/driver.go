package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//
//	DRIVERS - Conductores de ultima milla
//

type Driver struct {
	gorm.Model
	BusinessID     uint           `gorm:"not null;index;uniqueIndex:idx_driver_business_identification,priority:1"`
	FirstName      string         `gorm:"size:128;not null"`
	LastName       string         `gorm:"size:128;not null"`
	Email          string         `gorm:"size:255"`
	Phone          string         `gorm:"size:50;not null"`
	Identification string         `gorm:"size:50;not null;uniqueIndex:idx_driver_business_identification,priority:2"`
	Status         string         `gorm:"size:30;not null;default:'active';index"` // active, inactive, on_route
	PhotoURL       string         `gorm:"size:512"`
	LicenseType    string         `gorm:"size:20"` // A1, A2, B1, B2, C1
	LicenseExpiry  *time.Time
	WarehouseID    *uint          `gorm:"index"`
	Availability   datatypes.JSON `gorm:"type:jsonb"`
	Notes          *string        `gorm:"type:text"`

	// Relationships
	Business  Business   `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Warehouse *Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
