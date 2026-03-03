package entities

import "time"

// InventoryLevel representa el nivel de inventario de un producto en una bodega
type InventoryLevel struct {
	ID           uint
	ProductID    string
	WarehouseID  uint
	LocationID   *uint
	BusinessID   uint
	Quantity     int
	ReservedQty  int
	AvailableQty int
	MinStock     *int
	MaxStock     *int
	ReorderPoint *int
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Datos enriquecidos (no en DB)
	ProductName   string
	ProductSKU    string
	WarehouseName string
	WarehouseCode string
}
