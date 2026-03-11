package response

import (
	"encoding/json"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
)

type ProductResponse struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	ShortDescription string          `json:"short_description"`
	Price            float64         `json:"price"`
	CompareAtPrice   *float64        `json:"compare_at_price"`
	Currency         string          `json:"currency"`
	ImageURL         string          `json:"image_url"`
	Images           json.RawMessage `json:"images"`
	SKU              string          `json:"sku"`
	StockQuantity    int             `json:"stock_quantity"`
	Category         string          `json:"category"`
	Brand            string          `json:"brand"`
	IsFeatured       bool            `json:"is_featured"`
	CreatedAt        time.Time       `json:"created_at"`
}

type CatalogListResponse struct {
	Data       []ProductResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

func ProductFromEntity(p *entities.PublicProduct, imageURLBase string) ProductResponse {
	return ProductResponse{
		ID:               p.ID,
		Name:             p.Name,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Price:            p.Price,
		CompareAtPrice:   p.CompareAtPrice,
		Currency:         p.Currency,
		ImageURL:         buildFullImageURL(p.ImageURL, imageURLBase),
		Images:           p.Images,
		SKU:              p.SKU,
		StockQuantity:    p.StockQuantity,
		Category:         p.Category,
		Brand:            p.Brand,
		IsFeatured:       p.IsFeatured,
		CreatedAt:        p.CreatedAt,
	}
}

func ProductsFromEntities(products []entities.PublicProduct, imageURLBase string) []ProductResponse {
	result := make([]ProductResponse, len(products))
	for i := range products {
		result[i] = ProductFromEntity(&products[i], imageURLBase)
	}
	return result
}
