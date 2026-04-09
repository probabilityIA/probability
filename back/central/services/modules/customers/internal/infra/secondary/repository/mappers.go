package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func mapCustomerSummaryToEntity(m *models.CustomerSummary) *entities.CustomerSummary {
	return &entities.CustomerSummary{
		ID:                m.ID,
		CustomerID:        m.CustomerID,
		BusinessID:        m.BusinessID,
		TotalOrders:       m.TotalOrders,
		DeliveredOrders:   m.DeliveredOrders,
		CancelledOrders:   m.CancelledOrders,
		InProgressOrders:  m.InProgressOrders,
		TotalSpent:        m.TotalSpent,
		AvgTicket:         m.AvgTicket,
		TotalPaidOrders:   m.TotalPaidOrders,
		AvgDeliveryScore:  m.AvgDeliveryScore,
		FirstOrderAt:      m.FirstOrderAt,
		LastOrderAt:       m.LastOrderAt,
		PreferredPlatform: m.PreferredPlatform,
		LastUpdatedAt:     m.LastUpdatedAt,
	}
}

func mapCustomerSummaryFromEntity(e *entities.CustomerSummary) *models.CustomerSummary {
	m := &models.CustomerSummary{
		CustomerID:        e.CustomerID,
		BusinessID:        e.BusinessID,
		TotalOrders:       e.TotalOrders,
		DeliveredOrders:   e.DeliveredOrders,
		CancelledOrders:   e.CancelledOrders,
		InProgressOrders:  e.InProgressOrders,
		TotalSpent:        e.TotalSpent,
		AvgTicket:         e.AvgTicket,
		TotalPaidOrders:   e.TotalPaidOrders,
		AvgDeliveryScore:  e.AvgDeliveryScore,
		FirstOrderAt:      e.FirstOrderAt,
		LastOrderAt:       e.LastOrderAt,
		PreferredPlatform: e.PreferredPlatform,
		LastUpdatedAt:     e.LastUpdatedAt,
	}
	if e.ID > 0 {
		m.Model = gorm.Model{ID: e.ID}
	}
	return m
}

func mapCustomerAddressToEntity(m *models.CustomerAddress) entities.CustomerAddress {
	return entities.CustomerAddress{
		ID:         m.ID,
		CustomerID: m.CustomerID,
		BusinessID: m.BusinessID,
		Street:     m.Street,
		City:       m.City,
		State:      m.State,
		Country:    m.Country,
		PostalCode: m.PostalCode,
		TimesUsed:  m.TimesUsed,
		LastUsedAt: m.LastUsedAt,
	}
}

func mapCustomerAddressFromEntity(e *entities.CustomerAddress) *models.CustomerAddress {
	m := &models.CustomerAddress{
		CustomerID: e.CustomerID,
		BusinessID: e.BusinessID,
		Street:     e.Street,
		City:       e.City,
		State:      e.State,
		Country:    e.Country,
		PostalCode: e.PostalCode,
		TimesUsed:  e.TimesUsed,
		LastUsedAt: e.LastUsedAt,
	}
	if e.ID > 0 {
		m.Model = gorm.Model{ID: e.ID}
	}
	return m
}

func mapCustomerProductToEntity(m *models.CustomerProductHistory) entities.CustomerProductHistory {
	return entities.CustomerProductHistory{
		ID:             m.ID,
		CustomerID:     m.CustomerID,
		BusinessID:     m.BusinessID,
		ProductID:      m.ProductID,
		ProductName:    m.ProductName,
		ProductSKU:     m.ProductSKU,
		ProductImage:   m.ProductImage,
		TimesOrdered:   m.TimesOrdered,
		TotalQuantity:  m.TotalQuantity,
		TotalSpent:     m.TotalSpent,
		FirstOrderedAt: m.FirstOrderedAt,
		LastOrderedAt:  m.LastOrderedAt,
	}
}

func mapCustomerProductFromEntity(e *entities.CustomerProductHistory) *models.CustomerProductHistory {
	m := &models.CustomerProductHistory{
		CustomerID:     e.CustomerID,
		BusinessID:     e.BusinessID,
		ProductID:      e.ProductID,
		ProductName:    e.ProductName,
		ProductSKU:     e.ProductSKU,
		ProductImage:   e.ProductImage,
		TimesOrdered:   e.TimesOrdered,
		TotalQuantity:  e.TotalQuantity,
		TotalSpent:     e.TotalSpent,
		FirstOrderedAt: e.FirstOrderedAt,
		LastOrderedAt:  e.LastOrderedAt,
	}
	if e.ID > 0 {
		m.Model = gorm.Model{ID: e.ID}
	}
	return m
}

func mapCustomerOrderItemToEntity(m *models.CustomerOrderItem) entities.CustomerOrderItem {
	return entities.CustomerOrderItem{
		ID:           m.ID,
		CustomerID:   m.CustomerID,
		BusinessID:   m.BusinessID,
		OrderID:      m.OrderID,
		OrderNumber:  m.OrderNumber,
		ProductID:    m.ProductID,
		ProductName:  m.ProductName,
		ProductSKU:   m.ProductSKU,
		ProductImage: m.ProductImage,
		Quantity:     m.Quantity,
		UnitPrice:    m.UnitPrice,
		TotalPrice:   m.TotalPrice,
		OrderStatus:  m.OrderStatus,
		OrderedAt:    m.OrderedAt,
	}
}

func mapCustomerOrderItemFromEntity(e *entities.CustomerOrderItem) *models.CustomerOrderItem {
	m := &models.CustomerOrderItem{
		CustomerID:   e.CustomerID,
		BusinessID:   e.BusinessID,
		OrderID:      e.OrderID,
		OrderNumber:  e.OrderNumber,
		ProductID:    e.ProductID,
		ProductName:  e.ProductName,
		ProductSKU:   e.ProductSKU,
		ProductImage: e.ProductImage,
		Quantity:     e.Quantity,
		UnitPrice:    e.UnitPrice,
		TotalPrice:   e.TotalPrice,
		OrderStatus:  e.OrderStatus,
		OrderedAt:    e.OrderedAt,
	}
	if e.ID > 0 {
		m.Model = gorm.Model{ID: e.ID}
	}
	return m
}

func mapClientToEntity(m *models.Client) *entities.Client {
	return &entities.Client{
		ID:         m.ID,
		BusinessID: m.BusinessID,
		Name:       m.Name,
		Email:      m.Email,
		Phone:      m.Phone,
		Dni:        m.Dni,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
