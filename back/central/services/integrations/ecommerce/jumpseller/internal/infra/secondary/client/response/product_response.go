package response

import (
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

type ProductEnvelope struct {
	Product ProductResponse `json:"product"`
}

type ProductResponse struct {
	ID             int64                    `json:"id"`
	Name           string                   `json:"name"`
	SKU            string                   `json:"sku"`
	Barcode        string                   `json:"barcode"`
	Description    string                   `json:"description"`
	Price          float64                  `json:"price"`
	Stock          int                      `json:"stock"`
	StockUnlimited bool                     `json:"stock_unlimited"`
	Status         string                   `json:"status"`
	Weight         float64                  `json:"weight"`
	Height         float64                  `json:"height"`
	Width          float64                  `json:"width"`
	Length         float64                  `json:"length"`
	Diameter       float64                  `json:"diameter"`
	PackageFormat  string                   `json:"package_format"`
	Variants       []ProductVariantResponse `json:"variants"`
}

type ProductVariantResponse struct {
	ID             int64   `json:"id"`
	SKU            string  `json:"sku"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
	StockUnlimited bool    `json:"stock_unlimited"`
}

func (p ProductResponse) ToDomain() domain.JumpsellerProduct {
	variants := make([]domain.ProductVariant, 0, len(p.Variants))
	for _, v := range p.Variants {
		variants = append(variants, domain.ProductVariant{
			ID:             v.ID,
			SKU:            v.SKU,
			Price:          v.Price,
			Stock:          v.Stock,
			StockUnlimited: v.StockUnlimited,
		})
	}
	return domain.JumpsellerProduct{
		ID:             p.ID,
		Name:           p.Name,
		SKU:            p.SKU,
		Barcode:        p.Barcode,
		Description:    p.Description,
		Price:          p.Price,
		Stock:          p.Stock,
		StockUnlimited: p.StockUnlimited,
		Status:         p.Status,
		Weight:         p.Weight,
		Height:         p.Height,
		Width:          p.Width,
		Length:         p.Length,
		Diameter:       p.Diameter,
		PackageFormat:  p.PackageFormat,
		Variants:       variants,
	}
}

type UpdateProductStockRequest struct {
	Product UpdateStockFields `json:"product"`
}

type UpdateVariantStockRequest struct {
	Variant UpdateStockFields `json:"variant"`
}

type UpdateStockFields struct {
	Stock          int   `json:"stock"`
	StockUnlimited *bool `json:"stock_unlimited,omitempty"`
}

type CreateProductRequest struct {
	Product CreateProductFields `json:"product"`
}

type CreateProductFields struct {
	Name           string   `json:"name"`
	SKU            string   `json:"sku"`
	Price          float64  `json:"price"`
	Description    string   `json:"description,omitempty"`
	Stock          int      `json:"stock"`
	StockUnlimited bool     `json:"stock_unlimited"`
	Status         string   `json:"status"`
	Weight         *float64 `json:"weight,omitempty"`
	Height         *float64 `json:"height,omitempty"`
	Width          *float64 `json:"width,omitempty"`
	Length         *float64 `json:"length,omitempty"`
}

type UpdateProductRequest struct {
	Product UpdateProductFields `json:"product"`
}

type UpdateProductFields struct {
	Name        string   `json:"name,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Description string   `json:"description,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`
	Height      *float64 `json:"height,omitempty"`
	Width       *float64 `json:"width,omitempty"`
	Length      *float64 `json:"length,omitempty"`
}
