package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func ProductToEntity(m *models.Product) *entities.PublicProduct {
	return &entities.PublicProduct{
		ID:               m.ID,
		Name:             m.Name,
		Description:      m.Description,
		ShortDescription: m.ShortDescription,
		Price:            m.Price,
		CompareAtPrice:   m.CompareAtPrice,
		Currency:         m.Currency,
		ImageURL:         m.ImageURL,
		Images:           m.Images,
		SKU:              m.SKU,
		StockQuantity:    m.StockQuantity,
		Category:         m.Category,
		Brand:            m.Brand,
		IsFeatured:       m.IsFeatured,
		CreatedAt:        m.CreatedAt,
	}
}

func ProductsToEntities(products []models.Product) []entities.PublicProduct {
	result := make([]entities.PublicProduct, len(products))
	for i := range products {
		result[i] = *ProductToEntity(&products[i])
	}
	return result
}
