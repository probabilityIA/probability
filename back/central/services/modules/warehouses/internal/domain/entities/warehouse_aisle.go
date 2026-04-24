package entities

import "time"

type WarehouseAisle struct {
	ID         uint
	ZoneID     uint
	BusinessID uint
	Code       string
	Name       string
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
