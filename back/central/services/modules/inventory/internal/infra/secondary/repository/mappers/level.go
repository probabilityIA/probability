package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func LevelModelToEntity(m *models.InventoryLevel) *entities.InventoryLevel {
	return &entities.InventoryLevel{
		ID:           m.ID,
		ProductID:    m.ProductID,
		WarehouseID:  m.WarehouseID,
		LocationID:   m.LocationID,
		StateID:      m.StateID,
		BusinessID:   m.BusinessID,
		Quantity:     m.Quantity,
		ReservedQty:  m.ReservedQty,
		AvailableQty: m.AvailableQty,
		MinStock:     m.MinStock,
		MaxStock:     m.MaxStock,
		ReorderPoint: m.ReorderPoint,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func LevelEntityToModel(e *entities.InventoryLevel) *models.InventoryLevel {
	return &models.InventoryLevel{
		ProductID:    e.ProductID,
		WarehouseID:  e.WarehouseID,
		LocationID:   e.LocationID,
		BusinessID:   e.BusinessID,
		Quantity:     e.Quantity,
		ReservedQty:  e.ReservedQty,
		AvailableQty: e.AvailableQty,
		MinStock:     e.MinStock,
		MaxStock:     e.MaxStock,
		ReorderPoint: e.ReorderPoint,
	}
}
