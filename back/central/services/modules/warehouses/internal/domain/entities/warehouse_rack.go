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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
