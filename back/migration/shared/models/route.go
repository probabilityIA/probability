package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//
//	ROUTES - Rutas de ultima milla
//

type Route struct {
	gorm.Model
	BusinessID        uint       `gorm:"not null;index"`
	DriverID          *uint      `gorm:"index"`
	VehicleID         *uint      `gorm:"index"`
	Status            string     `gorm:"size:30;not null;default:'planned';index"` // planned, in_progress, completed, cancelled
	Date              time.Time  `gorm:"not null;index"`
	StartTime         *time.Time
	EndTime           *time.Time
	ActualStartTime   *time.Time
	ActualEndTime     *time.Time
	OriginWarehouseID *uint    `gorm:"index"`
	OriginAddress     string   `gorm:"size:500"`
	OriginLat         *float64 `gorm:"type:decimal(10,8)"`
	OriginLng         *float64 `gorm:"type:decimal(11,8)"`
	TotalStops        int      `gorm:"default:0"`
	CompletedStops    int      `gorm:"default:0"`
	FailedStops       int      `gorm:"default:0"`
	TotalDistanceKm   *float64 `gorm:"type:decimal(10,2)"`
	TotalDurationMin  *int
	OptimizedWaypoints datatypes.JSON `gorm:"type:jsonb"`
	Notes              *string        `gorm:"type:text"`

	// Relationships
	Business        Business    `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Driver          *Driver     `gorm:"foreignKey:DriverID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Vehicle         *Vehicle    `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	OriginWarehouse *Warehouse  `gorm:"foreignKey:OriginWarehouseID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Stops           []RouteStop `gorm:"foreignKey:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

//
//	ROUTE STOPS - Paradas de una ruta
//

type RouteStop struct {
	gorm.Model
	RouteID          uint       `gorm:"not null;index"`
	OrderID          *string    `gorm:"type:varchar(36);index"`
	Sequence         int        `gorm:"not null"`
	Status           string     `gorm:"size:30;not null;default:'pending';index"` // pending, arrived, delivered, failed, skipped
	Address          string     `gorm:"size:500"`
	City             string     `gorm:"size:128"`
	Lat              *float64   `gorm:"type:decimal(10,8)"`
	Lng              *float64   `gorm:"type:decimal(11,8)"`
	CustomerName     string     `gorm:"size:255"`
	CustomerPhone    string     `gorm:"size:50"`
	EstimatedArrival *time.Time
	ActualArrival    *time.Time
	ActualDeparture  *time.Time
	SignatureURL     string         `gorm:"size:512"`
	PhotoURL         string         `gorm:"size:512"`
	DeliveryNotes    *string        `gorm:"type:text"`
	FailureReason    *string        `gorm:"type:text"`
	Metadata         datatypes.JSON `gorm:"type:jsonb"`

	// Relationships
	Route Route  `gorm:"foreignKey:RouteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Order *Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
