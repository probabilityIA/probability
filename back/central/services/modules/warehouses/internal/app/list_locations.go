package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (uc *UseCase) ListLocations(ctx context.Context, params dtos.ListLocationsParams) ([]entities.WarehouseLocation, error) {
	// Verificar que la bodega pertenece al negocio
	_, err := uc.repo.GetByID(ctx, params.BusinessID, params.WarehouseID)
	if err != nil {
		return nil, err
	}

	return uc.repo.ListLocations(ctx, params)
}
