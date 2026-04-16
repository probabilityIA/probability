package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) UpdateOrderStatusMapping(ctx context.Context, id uint, mapping *entities.OrderStatusMapping) (*entities.OrderStatusMapping, error) {
	// Obtener el mapping actual
	current, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Actualizar solo los campos permitidos
	current.OriginalStatus = mapping.OriginalStatus
	current.OrderStatusID = mapping.OrderStatusID
	current.Description = mapping.Description

	if err := uc.repo.Update(ctx, current); err != nil {
		return nil, err
	}

	// Recargar para obtener relaciones actualizadas
	return uc.repo.GetByID(ctx, id)
}
