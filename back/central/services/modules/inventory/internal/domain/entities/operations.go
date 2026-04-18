package entities

import "time"

type PutawayRule struct {
	ID           uint
	BusinessID   uint
	ProductID    *string
	CategoryID   *uint
	TargetZoneID uint
	Priority     int
	Strategy     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PutawaySuggestion struct {
	ID                    uint
	BusinessID            uint
	ProductID             string
	RecommendedLocationID uint
	Quantity              int
	Status                string
	RuleID                *uint
	Reason                string
	ActualLocationID      *uint
	ConfirmedAt           *time.Time
	ConfirmedByID         *uint
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type ReplenishmentTask struct {
	ID             uint
	BusinessID     uint
	ProductID      string
	WarehouseID    uint
	FromLocationID *uint
	ToLocationID   *uint
	Quantity       int
	Status         string
	TriggeredBy    string
	AssignedToID   *uint
	AssignedAt     *time.Time
	CompletedAt    *time.Time
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CrossDockLink struct {
	ID                uint
	BusinessID        uint
	InboundShipmentID *uint
	OutboundOrderID   string
	ProductID         string
	Quantity          int
	Status            string
	ExecutedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ProductVelocity struct {
	ID          uint
	BusinessID  uint
	ProductID   string
	WarehouseID uint
	Period      string
	UnitsMoved  int
	Rank        string
	ComputedAt  time.Time
}
