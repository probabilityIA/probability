package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
	"time"
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
		Barcode:    p.Barcode,
		FamilyID:   p.FamilyID,

		// Información Básica
		Name:              p.Name,
		Title:             p.Title,
		Description:       p.Description,
		ShortDescription:  p.ShortDescription,
		Slug:              p.Slug,
		VariantLabel:      p.VariantLabel,
		VariantAttributes: p.VariantAttributes,
		VariantSignature:  p.VariantSignature,

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

// ToDBProductFamily convierte una familia de producto de dominio a modelo de base de datos.
func ToDBProductFamily(f *domain.ProductFamily) *models.ProductFamily {
	if f == nil {
		return nil
	}

	dbFamily := &models.ProductFamily{
		Model: gorm.Model{
			ID:        f.ID,
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		},
		BusinessID:  f.BusinessID,
		Name:        f.Name,
		Title:       f.Title,
		Description: f.Description,
		Slug:        f.Slug,
		Category:    f.Category,
		Brand:       f.Brand,
		ImageURL:    f.ImageURL,
		Status:      f.Status,
		IsActive:    f.IsActive,
		VariantAxes: f.VariantAxes,
		Metadata:    f.Metadata,
	}

	if f.DeletedAt != nil {
		dbFamily.DeletedAt = gorm.DeletedAt{Time: *f.DeletedAt, Valid: true}
	}

	return dbFamily
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
		Barcode:    p.Barcode,
		FamilyID:   p.FamilyID,

		// Información Básica
		Name:              p.Name,
		Title:             p.Title,
		Description:       p.Description,
		ShortDescription:  p.ShortDescription,
		Slug:              p.Slug,
		VariantLabel:      p.VariantLabel,
		VariantAttributes: p.VariantAttributes,
		VariantSignature:  p.VariantSignature,

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
		Family:   ToDomainProductFamily(p.Family),
	}
}

// ToDomainProductFamily convierte una familia de producto de base de datos a dominio.
func ToDomainProductFamily(f *models.ProductFamily) *domain.ProductFamily {
	if f == nil {
		return nil
	}

	var deletedAt *time.Time
	if f.DeletedAt.Valid {
		deletedAt = &f.DeletedAt.Time
	}

	return &domain.ProductFamily{
		ID:           f.ID,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
		DeletedAt:    deletedAt,
		BusinessID:   f.BusinessID,
		Name:         f.Name,
		Title:        f.Title,
		Description:  f.Description,
		Slug:         f.Slug,
		Category:     f.Category,
		Brand:        f.Brand,
		ImageURL:     f.ImageURL,
		Status:       f.Status,
		IsActive:     f.IsActive,
		VariantAxes:  f.VariantAxes,
		Metadata:     f.Metadata,
		VariantCount: f.VariantCount,
	}
}
