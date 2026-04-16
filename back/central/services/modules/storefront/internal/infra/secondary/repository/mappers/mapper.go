package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ProductToEntity maps a GORM product model to a storefront product entity
func ProductToEntity(m *models.Product) *entities.StorefrontProduct {
	return &entities.StorefrontProduct{
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
		Status:           m.Status,
		IsFeatured:       m.IsFeatured,
		CreatedAt:        m.CreatedAt,
	}
}

// ProductsToEntities maps a slice of GORM product models to entities
func ProductsToEntities(products []models.Product) []entities.StorefrontProduct {
	result := make([]entities.StorefrontProduct, len(products))
	for i := range products {
		result[i] = *ProductToEntity(&products[i])
	}
	return result
}

// OrderToEntity maps a GORM order model to a storefront order entity
func OrderToEntity(m *models.Order) *entities.StorefrontOrder {
	order := &entities.StorefrontOrder{
		ID:          m.ID,
		OrderNumber: m.OrderNumber,
		Status:      m.OriginalStatus,
		TotalAmount: m.TotalAmount,
		Currency:    m.Currency,
		CreatedAt:   m.CreatedAt,
	}

	for _, item := range m.OrderItems {
		productName := ""
		if item.Product.Name != "" {
			productName = item.Product.Name
		}
		var imageURL *string
		if item.Product.ImageURL != "" {
			url := item.Product.ImageURL
			imageURL = &url
		}

		order.Items = append(order.Items, entities.StorefrontOrderItem{
			ProductName: productName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
			ImageURL:    imageURL,
		})
	}

	return order
}

// OrdersToEntities maps a slice of GORM order models to entities
func OrdersToEntities(orders []models.Order) []entities.StorefrontOrder {
	result := make([]entities.StorefrontOrder, len(orders))
	for i := range orders {
		result[i] = *OrderToEntity(&orders[i])
	}
	return result
}
