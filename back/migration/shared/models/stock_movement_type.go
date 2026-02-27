package models

import "gorm.io/gorm"

// StockMovementType representa un tipo de movimiento de inventario (entrada, salida, ajuste, etc.)
type StockMovementType struct {
	gorm.Model
	Code        string `gorm:"size:50;uniqueIndex;not null"` // Código único del tipo: "inbound", "outbound", "adjustment", "transfer", "return", "sync"
	Name        string `gorm:"size:100;not null"`            // Nombre legible en español: "Entrada de mercancía", "Salida de mercancía", etc.
	Description string `gorm:"size:500"`                     // Descripción detallada del propósito de este tipo de movimiento
	IsActive    bool   `gorm:"default:true;index"`           // Indica si el tipo de movimiento está habilitado para su uso
	Direction   string `gorm:"size:10;not null"`             // Dirección del flujo de inventario: "in" (entrada), "out" (salida), "neutral" (sin cambio neto)
}

// TableName especifica el nombre de la tabla
func (StockMovementType) TableName() string {
	return "stock_movement_types"
}
