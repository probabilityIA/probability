package entities

import "time"

// StockMovementType representa un tipo de movimiento de inventario
type StockMovementType struct {
	ID          uint
	Code        string
	Name        string
	Description string
	IsActive    bool
	Direction   string // in, out, neutral
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
