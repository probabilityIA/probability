package entities

import "time"

type InventoryLevel struct {
	ID           uint
	ProductID    string
	WarehouseID  uint
	LocationID   *uint
	StateID      *uint
	BusinessID   uint
	Quantity     int
	ReservedQty  int
	AvailableQty int
	MinStock     *int
	MaxStock     *int
	ReorderPoint *int
	CreatedAt    time.Time
	UpdatedAt    time.Time

	ProductName      string
	ProductSKU       string
	WarehouseName    string
	WarehouseCode    string
	StateName        string
	LocationName     string
	LocationCode     string
}
