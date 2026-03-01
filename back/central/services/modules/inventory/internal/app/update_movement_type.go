package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) UpdateMovementType(ctx context.Context, dto dtos.UpdateStockMovementTypeDTO) (*entities.StockMovementType, error) {
	existing, err := uc.repo.GetMovementTypeByID(ctx, dto.ID)
	if err != nil {
		return nil, domainerrors.ErrMovementTypeNotFound
	}

	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.Description != "" {
		existing.Description = dto.Description
	}
	if dto.Direction != "" {
		existing.Direction = dto.Direction
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}

	if err := uc.repo.UpdateMovementType(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}
