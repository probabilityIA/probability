package entities

import "time"

// Product representa un producto en el dominio
// ✅ ENTIDAD PURA - SIN TAGS
type Product struct {
	ID         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	BusinessID uint
	SKU        string
	Name       string
	ExternalID string
}

// ToDomainProduct convierte un modelo de BD a dominio
func ToDomainProduct(p interface{}) *Product {
	// Nota: Esto se implementará correctamente en el mapper,
	// aquí solo definimos la estructura.
	return &Product{}
}
