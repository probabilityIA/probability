package entities

import "time"

type WarehouseRackLevel struct {
	ID         uint
	RackID     uint
	BusinessID uint
	Code       string
	Ordinal    int
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
