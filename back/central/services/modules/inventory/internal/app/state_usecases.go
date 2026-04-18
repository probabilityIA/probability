package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) ListInventoryStates(ctx context.Context) ([]entities.InventoryState, error) {
	return uc.repo.ListInventoryStates(ctx)
}

func (uc *useCase) ChangeInventoryState(ctx context.Context, dto request.ChangeInventoryStateDTO) (*entities.StockMovement, error) {
	if dto.Quantity <= 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}
	if dto.FromStateCode == dto.ToStateCode {
		return nil, domainerrors.ErrStateTransition
	}

	movement, err := uc.repo.ChangeStateTx(ctx, dtos.ChangeInventoryStateTxParams{
		LevelID:       dto.LevelID,
		FromStateCode: dto.FromStateCode,
		ToStateCode:   dto.ToStateCode,
		Quantity:      dto.Quantity,
		Reason:        dto.Reason,
		BusinessID:    dto.BusinessID,
		CreatedByID:   dto.CreatedByID,
	})
	if err != nil {
		return nil, err
	}
	return movement, nil
}
