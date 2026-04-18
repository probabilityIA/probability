package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func MovementTypeModelToEntity(m *models.StockMovementType) *entities.StockMovementType {
	return &entities.StockMovementType{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		IsActive:    m.IsActive,
		Direction:   m.Direction,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func MovementTypeEntityToModel(e *entities.StockMovementType) *models.StockMovementType {
	return &models.StockMovementType{
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		IsActive:    e.IsActive,
		Direction:   e.Direction,
	}
}
