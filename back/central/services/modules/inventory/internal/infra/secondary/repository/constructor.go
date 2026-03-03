package repository

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// Repository implementa ports.IRepository
type Repository struct {
	db    db.IDatabase
	cache IInventoryCache
}

// New crea una nueva instancia del repositorio
func New(database db.IDatabase, cache IInventoryCache) ports.IRepository {
	return &Repository{db: database, cache: cache}
}

func inventoryLevelModelToEntity(m *models.InventoryLevel) *entities.InventoryLevel {
	return &entities.InventoryLevel{
		ID:           m.ID,
		ProductID:    m.ProductID,
		WarehouseID:  m.WarehouseID,
		LocationID:   m.LocationID,
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

func stockMovementModelToEntity(m *models.StockMovement) *entities.StockMovement {
	return &entities.StockMovement{
		ID:             m.ID,
		ProductID:      m.ProductID,
		WarehouseID:    m.WarehouseID,
		LocationID:     m.LocationID,
		BusinessID:     m.BusinessID,
		MovementTypeID: m.MovementTypeID,
		Reason:         m.Reason,
		Quantity:       m.Quantity,
		PreviousQty:    m.PreviousQty,
		NewQty:         m.NewQty,
		ReferenceType:  m.ReferenceType,
		ReferenceID:    m.ReferenceID,
		IntegrationID:  m.IntegrationID,
		Notes:          m.Notes,
		CreatedByID:    m.CreatedByID,
		CreatedAt:      m.CreatedAt,
	}
}
