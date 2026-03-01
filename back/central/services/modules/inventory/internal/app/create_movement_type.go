package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (uc *useCase) CreateMovementType(ctx context.Context, dto dtos.CreateStockMovementTypeDTO) (*entities.StockMovementType, error) {
	movType := &entities.StockMovementType{
		Code:        dto.Code,
		Name:        dto.Name,
		Description: dto.Description,
		Direction:   dto.Direction,
		IsActive:    true,
	}
	return uc.repo.CreateMovementType(ctx, movType)
}
