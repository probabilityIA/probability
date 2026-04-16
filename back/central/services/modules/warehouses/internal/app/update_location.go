package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (uc *UseCase) UpdateLocation(ctx context.Context, dto dtos.UpdateLocationDTO) (*entities.WarehouseLocation, error) {
	// Verificar que la bodega existe y pertenece al negocio
	_, err := uc.repo.GetByID(ctx, dto.BusinessID, dto.WarehouseID)
	if err != nil {
		return nil, err
	}

	// Verificar que la ubicación existe
	_, err = uc.repo.GetLocationByID(ctx, dto.WarehouseID, dto.ID)
	if err != nil {
		return nil, err
	}

	// Verificar código duplicado (excluyendo la actual)
	exists, err := uc.repo.LocationExistsByCode(ctx, dto.WarehouseID, dto.Code, &dto.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateLocCode
	}

	location := &entities.WarehouseLocation{
		ID:            dto.ID,
		WarehouseID:   dto.WarehouseID,
		Name:          dto.Name,
		Code:          dto.Code,
		Type:          dto.Type,
		IsActive:      dto.IsActive,
		IsFulfillment: dto.IsFulfillment,
		Capacity:      dto.Capacity,
	}

	return uc.repo.UpdateLocation(ctx, location)
}
