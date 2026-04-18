package models

import (
	"time"

	"gorm.io/gorm"
)

type CycleCountPlan struct {
	gorm.Model
	BusinessID     uint   `gorm:"not null;index"`
	WarehouseID    uint   `gorm:"not null;index"`
	Name           string `gorm:"size:255;not null"`
	Strategy       string `gorm:"size:20;default:'abc';index"`
	FrequencyDays  int    `gorm:"default:30"`
	NextRunAt      *time.Time
	IsActive       bool `gorm:"default:true;index"`

	Business  Business  `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Warehouse Warehouse `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CycleCountPlan) TableName() string {
	return "cycle_count_plans"
}

type CycleCountTask struct {
	gorm.Model
	PlanID       uint   `gorm:"not null;index"`
	BusinessID   uint   `gorm:"not null;index"`
	WarehouseID  uint   `gorm:"not null;index"`
	ScopeType    string `gorm:"size:20;default:'zone'"`
	ScopeID      *uint  `gorm:"index"`
	Status       string `gorm:"size:20;default:'pending';index"`
	AssignedToID *uint  `gorm:"index"`
	StartedAt    *time.Time
	FinishedAt   *time.Time

	Plan      CycleCountPlan `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business  Business       `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Warehouse Warehouse      `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CycleCountTask) TableName() string {
	return "cycle_count_tasks"
}

type CycleCountLine struct {
	gorm.Model
	TaskID      uint   `gorm:"not null;index"`
	BusinessID  uint   `gorm:"not null;index"`
	ProductID   string `gorm:"type:varchar(64);not null;index"`
	LocationID  *uint  `gorm:"index"`
	LotID       *uint  `gorm:"index"`
	ExpectedQty int    `gorm:"not null"`
	CountedQty  *int
	Variance    int    `gorm:"default:0"`
	Status      string `gorm:"size:20;default:'pending';index"`

	Task     CycleCountTask `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business       `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CycleCountLine) TableName() string {
	return "cycle_count_lines"
}

type InventoryDiscrepancy struct {
	gorm.Model
	TaskID               uint   `gorm:"not null;index"`
	LineID               uint   `gorm:"not null;index"`
	BusinessID           uint   `gorm:"not null;index"`
	Status               string `gorm:"size:20;default:'open';index"`
	ResolutionMovementID *uint  `gorm:"index"`
	ReviewedByID         *uint  `gorm:"index"`
	ReviewedAt           *time.Time
	Notes                string `gorm:"type:text"`

	Task     CycleCountTask `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Line     CycleCountLine `gorm:"foreignKey:LineID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business Business       `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (InventoryDiscrepancy) TableName() string {
	return "inventory_discrepancies"
}
