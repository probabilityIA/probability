package entities

import "time"

type WarehouseZone struct {
	ID          uint
	WarehouseID uint
	BusinessID  uint
	Code        string
	Name        string
	Purpose     string
	IsActive    bool
	ColorHex    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
