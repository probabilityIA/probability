package domain

import (
	"time"
)

// ───────────────────────────────────────────
//
//	PRODUCT DTOs
//
// ───────────────────────────────────────────

// CreateProductRequest representa la solicitud para crear un producto
type CreateProductRequest struct {
	BusinessID uint   `json:"business_id" binding:"required"`
	SKU        string `json:"sku" binding:"required,max=128"`
	Name       string `json:"name" binding:"required,max=255"`
	ExternalID string `json:"external_id" binding:"max=255"`
}

// UpdateProductRequest representa la solicitud para actualizar un producto
type UpdateProductRequest struct {
	SKU        *string `json:"sku" binding:"omitempty,max=128"`
	Name       *string `json:"name" binding:"omitempty,max=255"`
	ExternalID *string `json:"external_id" binding:"omitempty,max=255"`
}

// ProductResponse representa la respuesta de un producto
type ProductResponse struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	BusinessID uint       `json:"business_id"`
	SKU        string     `json:"sku"`
	Name       string     `json:"name"`
	ExternalID string     `json:"external_id"`
}

// ProductsListResponse representa la respuesta paginada de productos
type ProductsListResponse struct {
	Data       []ProductResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

