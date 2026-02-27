package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (uc *UseCase) UpdateWarehouse(ctx context.Context, dto dtos.UpdateWarehouseDTO) (*entities.Warehouse, error) {
	// Verificar que existe
	_, err := uc.repo.GetByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	// Verificar código duplicado (excluyendo la actual)
	exists, err := uc.repo.ExistsByCode(ctx, dto.BusinessID, dto.Code, &dto.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateCode
	}

	// Si se marca como default, quitar default de las demás
	if dto.IsDefault {
		if err := uc.repo.ClearDefault(ctx, dto.BusinessID, dto.ID); err != nil {
			return nil, err
		}
	}

	warehouse := &entities.Warehouse{
		ID:            dto.ID,
		BusinessID:    dto.BusinessID,
		Name:          dto.Name,
		Code:          dto.Code,
		Address:       dto.Address,
		City:          dto.City,
		State:         dto.State,
		Country:       dto.Country,
		ZipCode:       dto.ZipCode,
		Phone:         dto.Phone,
		ContactName:   dto.ContactName,
		ContactEmail:  dto.ContactEmail,
		IsActive:      dto.IsActive,
		IsDefault:     dto.IsDefault,
		IsFulfillment: dto.IsFulfillment,
	}

	return uc.repo.Update(ctx, warehouse)
}
