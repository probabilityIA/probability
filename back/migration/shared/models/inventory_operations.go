package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PutawayRule struct {
	gorm.Model
	BusinessID   uint    `gorm:"not null;index"`
	ProductID    *string `gorm:"type:varchar(64);index"`
	CategoryID   *uint   `gorm:"index"`
	TargetZoneID uint    `gorm:"not null;index"`
	Priority     int     `gorm:"default:0;index"`
	Strategy     string  `gorm:"size:30;default:'nearest_empty'"`
	IsActive     bool    `gorm:"default:true;index"`

	Business Business      `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Zone     WarehouseZone `gorm:"foreignKey:TargetZoneID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (PutawayRule) TableName() string {
	return "putaway_rules"
}

type PutawaySuggestion struct {
	gorm.Model
	BusinessID           uint   `gorm:"not null;index"`
	ProductID            string `gorm:"type:varchar(64);not null;index"`
	RecommendedLocationID uint  `gorm:"not null;index"`
	Quantity             int    `gorm:"not null"`
	Status               string `gorm:"size:20;default:'pending';index"`
	RuleID               *uint  `gorm:"index"`
	Reason               string `gorm:"size:255"`
	ActualLocationID     *uint  `gorm:"index"`
	ConfirmedAt          *time.Time
	ConfirmedByID        *uint

	Business            Business          `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RecommendedLocation WarehouseLocation `gorm:"foreignKey:RecommendedLocationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (PutawaySuggestion) TableName() string {
	return "putaway_suggestions"
}

type ReplenishmentTask struct {
	gorm.Model
	BusinessID     uint   `gorm:"not null;index"`
	ProductID      string `gorm:"type:varchar(64);not null;index"`
	WarehouseID    uint   `gorm:"not null;index"`
	FromLocationID *uint  `gorm:"index"`
	ToLocationID   *uint  `gorm:"index"`
	Quantity       int    `gorm:"not null"`
	Status         string `gorm:"size:20;default:'pending';index"`
	TriggeredBy    string `gorm:"size:10;default:'auto'"`
	AssignedToID   *uint  `gorm:"index"`
	AssignedAt     *time.Time
	CompletedAt    *time.Time
	Notes          string `gorm:"type:text"`

	Business  Business  `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Warehouse Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ReplenishmentTask) TableName() string {
	return "replenishment_tasks"
}

type CrossDockLink struct {
	gorm.Model
	BusinessID         uint   `gorm:"not null;index"`
	InboundShipmentID  *uint  `gorm:"index"`
	OutboundOrderID    string `gorm:"type:varchar(64);not null;index"`
	ProductID          string `gorm:"type:varchar(64);not null;index"`
	Quantity           int    `gorm:"not null"`
	Status             string `gorm:"size:20;default:'pending';index"`
	ExecutedAt         *time.Time
	Metadata           datatypes.JSON `gorm:"type:jsonb"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CrossDockLink) TableName() string {
	return "cross_dock_links"
}

type ProductVelocity struct {
	gorm.Model
	BusinessID  uint   `gorm:"not null;index;uniqueIndex:idx_velocity_scope,priority:1"`
	ProductID   string `gorm:"type:varchar(64);not null;uniqueIndex:idx_velocity_scope,priority:2"`
	WarehouseID uint   `gorm:"not null;uniqueIndex:idx_velocity_scope,priority:3"`
	Period      string `gorm:"size:10;not null;uniqueIndex:idx_velocity_scope,priority:4"`
	UnitsMoved  int    `gorm:"not null;default:0"`
	Rank        string `gorm:"size:1;default:'C';index"`
	ComputedAt  time.Time

	Business  Business  `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Warehouse Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (ProductVelocity) TableName() string {
	return "product_velocities"
}
