package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	VEHICLES - Vehiculos de ultima milla
//
// ───────────────────────────────────────────

type Vehicle struct {
	gorm.Model
	BusinessID         uint   `gorm:"not null;index;uniqueIndex:idx_vehicle_business_plate,priority:1"`
	Type               string `gorm:"size:30;not null;index"` // motorcycle, car, van, truck
	LicensePlate       string `gorm:"size:20;not null;uniqueIndex:idx_vehicle_business_plate,priority:2"`
	Brand              string `gorm:"size:64"`
	VehicleModel       string `gorm:"size:64"`
	Year               *int
	Color              string         `gorm:"size:30"`
	Status             string         `gorm:"size:30;not null;default:'active';index"` // active, inactive, in_maintenance
	WeightCapacityKg   *float64       `gorm:"type:decimal(10,2)"`
	VolumeCapacityM3   *float64       `gorm:"type:decimal(10,2)"`
	PhotoURL           string         `gorm:"size:512"`
	InsuranceExpiry    *time.Time
	RegistrationExpiry *time.Time
	Metadata           datatypes.JSON `gorm:"type:jsonb"`

	// Relationships
	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
