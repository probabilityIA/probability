package entities

import "time"

// Client representa un cliente en el dominio
type Client struct {
	ID         uint
	BusinessID uint
	Name       string
	Email      string
	Phone      string
	Dni        *string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Stats (solo disponibles en GetClient - vienen de la tabla orders)
	OrderCount  int64
	TotalSpent  float64
	LastOrderAt *time.Time
}
