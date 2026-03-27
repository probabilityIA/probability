package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

// QuantityDiscountResponse respuesta de descuento por cantidad
type QuantityDiscountResponse struct {
	ID              uint      `json:"id"`
	BusinessID      uint      `json:"business_id"`
	ProductID       *string   `json:"product_id"`
	ProductName     string    `json:"product_name"`
	MinQuantity     int       `json:"min_quantity"`
	DiscountPercent float64   `json:"discount_percent"`
	IsActive        bool      `json:"is_active"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// FromDiscountEntity convierte una entidad de dominio a response
func FromDiscountEntity(d *entities.QuantityDiscount) QuantityDiscountResponse {
	return QuantityDiscountResponse{
		ID:              d.ID,
		BusinessID:      d.BusinessID,
		ProductID:       d.ProductID,
		ProductName:     d.ProductName,
		MinQuantity:     d.MinQuantity,
		DiscountPercent: d.DiscountPercent,
		IsActive:        d.IsActive,
		Description:     d.Description,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

// QuantityDiscountsListResponse respuesta paginada de descuentos
type QuantityDiscountsListResponse struct {
	Data       []QuantityDiscountResponse `json:"data"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	PageSize   int                        `json:"page_size"`
	TotalPages int                        `json:"total_pages"`
}
