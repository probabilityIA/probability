package entities

import "time"

// StockMovement representa un movimiento de inventario
type StockMovement struct {
	ID             uint
	ProductID      string
	WarehouseID    uint
	LocationID     *uint
	BusinessID     uint
	MovementTypeID uint   // FK a stock_movement_types
	Reason         string
	Quantity       int // positivo=entrada, negativo=salida
	PreviousQty    int
	NewQty         int
	ReferenceType  *string // order, shipment, manual, sync
	ReferenceID    *string
	IntegrationID  *uint
	Notes          string
	CreatedByID    *uint
	CreatedAt      time.Time

	// Datos enriquecidos (no en DB)
	MovementTypeCode string
	MovementTypeName string
	ProductName      string
	ProductSKU       string
	WarehouseName    string
}
