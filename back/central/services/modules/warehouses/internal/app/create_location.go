package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (uc *UseCase) CreateLocation(ctx context.Context, dto dtos.CreateLocationDTO) (*entities.WarehouseLocation, error) {
	// Verificar que la bodega existe y pertenece al negocio
	_, err := uc.repo.GetByID(ctx, dto.BusinessID, dto.WarehouseID)
	if err != nil {
		return nil, err
	}

	// Verificar código duplicado
	exists, err := uc.repo.LocationExistsByCode(ctx, dto.WarehouseID, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateLocCode
	}

	location := &entities.WarehouseLocation{
		WarehouseID:   dto.WarehouseID,
		Name:          dto.Name,
		Code:          dto.Code,
		Type:          dto.Type,
		IsActive:      dto.IsActive,
		IsFulfillment: dto.IsFulfillment,
		Capacity:      dto.Capacity,
	}

	return uc.repo.CreateLocation(ctx, location)
}
