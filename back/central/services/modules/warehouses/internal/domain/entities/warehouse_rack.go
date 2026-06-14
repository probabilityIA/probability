package entities

import "time"

type WarehouseRack struct {
	ID          uint
	AisleID     uint
	BusinessID  uint
	Code        string
	Name        string
	LevelsCount int
	IsActive    bool
	WidthCm     float64
	DepthCm     float64
	HeightCm    float64
	Side        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
