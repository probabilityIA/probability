package models

import "gorm.io/gorm"

// InventoryLevel representa el nivel de inventario de un producto en una bodega
type InventoryLevel struct {
	gorm.Model
	ProductID    string `gorm:"type:varchar(64);not null;uniqueIndex:idx_inventory_product_warehouse,priority:1"` // ID del producto (FK a products)
	WarehouseID  uint   `gorm:"not null;uniqueIndex:idx_inventory_product_warehouse,priority:2"`                  // ID de la bodega donde se almacena
	LocationID   *uint  `gorm:"index"`                                                                            // ID de la ubicación dentro de la bodega (nil = sin ubicación específica)
	BusinessID   uint   `gorm:"not null;index"`                                                                   // ID del negocio propietario
	Quantity     int    `gorm:"default:0;not null"`                                                               // Cantidad total de stock en esta bodega
	ReservedQty  int    `gorm:"default:0;not null"`                                                               // Cantidad reservada/comprometida por órdenes pendientes
	AvailableQty int    `gorm:"default:0;not null"`                                                               // Cantidad disponible para venta (Quantity - ReservedQty)
	MinStock     *int   //                                                                                          Nivel mínimo de stock para alertas (nil = sin mínimo)
	MaxStock     *int   //                                                                                          Nivel máximo de stock para control (nil = sin máximo)
	ReorderPoint *int   //                                                                                          Punto de reorden: nivel en el que se debe reabastecer (nil = sin punto de reorden)

	// Relaciones
	Product   Product            `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`   // Producto al que pertenece este nivel
	Warehouse Warehouse          `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Bodega donde se almacena
	Location  *WarehouseLocation `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Ubicación específica dentro de la bodega
	Business  Business           `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`  // Negocio propietario del inventario
}

// TableName especifica el nombre de la tabla
func (InventoryLevel) TableName() string {
	return "inventory_levels"
}

// StockMovement representa un movimiento de inventario (entrada, salida, ajuste, transferencia)
type StockMovement struct {
	gorm.Model
	ProductID      string  `gorm:"type:varchar(64);not null;index"`                                                                    // ID del producto afectado por el movimiento
	WarehouseID    uint    `gorm:"not null;index"`                                                                                      // ID de la bodega donde ocurre el movimiento
	LocationID     *uint   `gorm:"index"`                                                                                               // ID de la ubicación dentro de la bodega (nil = movimiento general)
	BusinessID     uint    `gorm:"not null;index"`                                                                                      // ID del negocio propietario
	MovementTypeID uint    `gorm:"not null;index"`                                                                                      // FK al tipo de movimiento (stock_movement_types)
	Reason         string  `gorm:"size:255"`                                                                                            // Motivo o justificación del movimiento
	Quantity       int     `gorm:"not null"`                                                                                            // Cantidad movida: positivo=entrada, negativo=salida
	PreviousQty    int     `gorm:"not null"`                                                                                            // Cantidad de stock antes del movimiento
	NewQty         int     `gorm:"not null"`                                                                                            // Cantidad de stock después del movimiento
	ReferenceType  *string `gorm:"size:50"`                                                                                             // Tipo de referencia del origen: order, shipment, manual, sync
	ReferenceID    *string `gorm:"size:64"`                                                                                             // ID de la referencia de origen (ej: ID de la orden)
	IntegrationID  *uint   `gorm:"index"`                                                                                               // ID de la integración que originó el movimiento (nil = manual)
	Notes          string  `gorm:"type:text"`                                                                                           // Notas o comentarios adicionales sobre el movimiento
	CreatedByID    *uint   `gorm:"index"`                                                                                               // ID del usuario que creó el movimiento (nil = sistema)

	// Relaciones
	MovementType StockMovementType `gorm:"foreignKey:MovementTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // Tipo de movimiento (entrada, salida, ajuste, etc.)
	Product      Product           `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`        // Producto afectado
	Warehouse    Warehouse         `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`      // Bodega donde ocurre el movimiento
	Business     Business          `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`       // Negocio propietario
	CreatedBy    *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`     // Usuario que creó el movimiento
}

// TableName especifica el nombre de la tabla
func (StockMovement) TableName() string {
	return "stock_movements"
}
