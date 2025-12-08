package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ToDBProduct convierte un producto de dominio a modelo de base de datos
func ToDBProduct(p *domain.Product) *models.Product {
	if p == nil {
		return nil
	}
	return &models.Product{
		// Timestamps
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,

		// Identificadores
		BusinessID: p.BusinessID,
		SKU:        p.SKU,
		ExternalID: p.ExternalID,

		// Información Básica
		Name:             p.Name,
		Title:            p.Title,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Slug:             p.Slug,

		// Pricing
		Price:          p.Price,
		CompareAtPrice: p.CompareAtPrice,
		CostPrice:      p.CostPrice,
		Currency:       p.Currency,

		// Inventory
		StockQuantity:     p.StockQuantity,
		TrackInventory:    p.TrackInventory,
		AllowBackorder:    p.AllowBackorder,
		LowStockThreshold: p.LowStockThreshold,

		// Media
		ImageURL: p.ImageURL,
		Images:   p.Images,
		VideoURL: p.VideoURL,

		// Dimensiones y Peso
		Weight:        p.Weight,
		WeightUnit:    p.WeightUnit,
		Length:        p.Length,
		Width:         p.Width,
		Height:        p.Height,
		DimensionUnit: p.DimensionUnit,

		// Categorización
		Category: p.Category,
		Tags:     p.Tags,
		Brand:    p.Brand,

		// Estado
		Status:     p.Status,
		IsActive:   p.IsActive,
		IsFeatured: p.IsFeatured,

		// Metadata
		Metadata: p.Metadata,
	}
}

// ToDomainProduct convierte un producto de base de datos a dominio
func ToDomainProduct(p *models.Product) *domain.Product {
	if p == nil {
		return nil
	}
	return &domain.Product{
		// Timestamps
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,

		// Identificadores
		BusinessID: p.BusinessID,
		SKU:        p.SKU,
		ExternalID: p.ExternalID,

		// Información Básica
		Name:             p.Name,
		Title:            p.Title,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Slug:             p.Slug,

		// Pricing
		Price:          p.Price,
		CompareAtPrice: p.CompareAtPrice,
		CostPrice:      p.CostPrice,
		Currency:       p.Currency,

		// Inventory
		StockQuantity:     p.StockQuantity,
		TrackInventory:    p.TrackInventory,
		AllowBackorder:    p.AllowBackorder,
		LowStockThreshold: p.LowStockThreshold,

		// Media
		ImageURL: p.ImageURL,
		Images:   p.Images,
		VideoURL: p.VideoURL,

		// Dimensiones y Peso
		Weight:        p.Weight,
		WeightUnit:    p.WeightUnit,
		Length:        p.Length,
		Width:         p.Width,
		Height:        p.Height,
		DimensionUnit: p.DimensionUnit,

		// Categorización
		Category: p.Category,
		Tags:     p.Tags,
		Brand:    p.Brand,

		// Estado
		Status:     p.Status,
		IsActive:   p.IsActive,
		IsFeatured: p.IsFeatured,

		// Metadata
		Metadata: p.Metadata,
	}
}
