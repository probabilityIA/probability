package entities

import "time"

type CycleCountPlan struct {
	ID            uint
	BusinessID    uint
	WarehouseID   uint
	Name          string
	Strategy      string
	FrequencyDays int
	NextRunAt     *time.Time
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CycleCountTask struct {
	ID           uint
	PlanID       uint
	BusinessID   uint
	WarehouseID  uint
	ScopeType    string
	ScopeID      *uint
	Status       string
	AssignedToID *uint
	StartedAt    *time.Time
	FinishedAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CycleCountLine struct {
	ID          uint
	TaskID      uint
	BusinessID  uint
	ProductID   string
	LocationID  *uint
	LotID       *uint
	ExpectedQty int
	CountedQty  *int
	Variance    int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InventoryDiscrepancy struct {
	ID                   uint
	TaskID               uint
	LineID               uint
	BusinessID           uint
	Status               string
	ResolutionMovementID *uint
	ReviewedByID         *uint
	ReviewedAt           *time.Time
	Notes                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type KardexEntry struct {
	MovementID       uint
	CreatedAt        time.Time
	MovementTypeCode string
	MovementTypeName string
	Quantity         int
	PreviousQty      int
	NewQty           int
	RunningBalance   int
	Reason           string
	ReferenceType    *string
	ReferenceID      *string
	LocationID       *uint
	LotID            *uint
}
