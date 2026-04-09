package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

type CustomerProductResponse struct {
	ID             uint      `json:"id"`
	CustomerID     uint      `json:"customer_id"`
	BusinessID     uint      `json:"business_id"`
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ProductSKU     string    `json:"product_sku"`
	ProductImage   *string   `json:"product_image"`
	TimesOrdered   int       `json:"times_ordered"`
	TotalQuantity  int       `json:"total_quantity"`
	TotalSpent     float64   `json:"total_spent"`
	FirstOrderedAt time.Time `json:"first_ordered_at"`
	LastOrderedAt  time.Time `json:"last_ordered_at"`
}

func ProductFromEntity(p *entities.CustomerProductHistory) CustomerProductResponse {
	return CustomerProductResponse{
		ID:             p.ID,
		CustomerID:     p.CustomerID,
		BusinessID:     p.BusinessID,
		ProductID:      p.ProductID,
		ProductName:    p.ProductName,
		ProductSKU:     p.ProductSKU,
		ProductImage:   p.ProductImage,
		TimesOrdered:   p.TimesOrdered,
		TotalQuantity:  p.TotalQuantity,
		TotalSpent:     p.TotalSpent,
		FirstOrderedAt: p.FirstOrderedAt,
		LastOrderedAt:  p.LastOrderedAt,
	}
}

type CustomerProductListResponse struct {
	Data       []CustomerProductResponse `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	PageSize   int                       `json:"page_size"`
	TotalPages int                       `json:"total_pages"`
}
