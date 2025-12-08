package domain

import (
	"time"
)

// Product representa un producto en el dominio
type Product struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	BusinessID uint       `json:"business_id"`
	SKU        string     `json:"sku"`
	Name       string     `json:"name"`
	ExternalID string     `json:"external_id"`
}

// Client representa un cliente en el dominio
type Client struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	BusinessID uint       `json:"business_id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Dni        *string    `json:"dni"`
}

// ToDomainProduct convierte un modelo de BD a dominio
func ToDomainProduct(p interface{}) *Product {
	// Nota: Esto se implementará correctamente en el mapper,
	// aquí solo definimos la estructura.
	return &Product{}
}
