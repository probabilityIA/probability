package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func MovementModelToEntity(m *models.StockMovement) *entities.StockMovement {
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

func MovementEntityToModel(e *entities.StockMovement) *models.StockMovement {
	return &models.StockMovement{
		ProductID:      e.ProductID,
		WarehouseID:    e.WarehouseID,
		LocationID:     e.LocationID,
		BusinessID:     e.BusinessID,
		MovementTypeID: e.MovementTypeID,
		Reason:         e.Reason,
		Quantity:       e.Quantity,
		PreviousQty:    e.PreviousQty,
		NewQty:         e.NewQty,
		ReferenceType:  e.ReferenceType,
		ReferenceID:    e.ReferenceID,
		IntegrationID:  e.IntegrationID,
		Notes:          e.Notes,
		CreatedByID:    e.CreatedByID,
	}
}
