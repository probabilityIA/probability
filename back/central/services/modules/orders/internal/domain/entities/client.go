package entities

import "time"

// Client representa un cliente en el dominio
// ✅ ENTIDAD PURA - SIN TAGS
type Client struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	BusinessID uint
	Name       string
	Email      *string
	Phone      string
	Dni        *string
}
