package entities

import "time"

// WarehouseLocation representa una ubicación dentro de una bodega
type WarehouseLocation struct {
	ID            uint
	WarehouseID   uint
	Name          string
	Code          string
	Type          string // storage, picking, packing, receiving, shipping
	IsActive      bool
	IsFulfillment bool
	Capacity      *int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
