package entities

import "time"

type WarehouseLocationFlags struct {
	IsPicking    bool `json:"is_picking"`
	IsBulk       bool `json:"is_bulk"`
	IsQuarantine bool `json:"is_quarantine"`
	IsDamaged    bool `json:"is_damaged"`
	IsReturns    bool `json:"is_returns"`
	IsCrossDock  bool `json:"is_cross_dock"`
	IsHazmat     bool `json:"is_hazmat"`
}

type WarehouseLocation struct {
	ID            uint
	WarehouseID   uint
	LevelID       *uint
	Name          string
	Code          string
	Type          string
	IsActive      bool
	IsFulfillment bool
	Capacity      *int
	MaxWeightKg   *float64
	MaxVolumeCm3  *float64
	LengthCm      *float64
	WidthCm       *float64
	HeightCm      *float64
	Priority      int
	Flags         WarehouseLocationFlags
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
