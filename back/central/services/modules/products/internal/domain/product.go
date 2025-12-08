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
	BusinessID uint        `json:"business_id"`
	SKU        string      `json:"sku"`
	Name       string      `json:"name"`
	ExternalID string      `json:"external_id"`
}

