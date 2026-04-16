package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func (uc *UseCase) CreateWarehouse(ctx context.Context, dto dtos.CreateWarehouseDTO) (*entities.Warehouse, error) {
	// Verificar código duplicado
	exists, err := uc.repo.ExistsByCode(ctx, dto.BusinessID, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrDuplicateCode
	}

	// Si es default, quitar default de las demás
	if dto.IsDefault {
		if err := uc.repo.ClearDefault(ctx, dto.BusinessID, 0); err != nil {
			return nil, err
		}
	}

	warehouse := &entities.Warehouse{
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
		Company:       dto.Company,
		FirstName:     dto.FirstName,
		LastName:      dto.LastName,
		Email:         dto.Email,
		Suburb:        dto.Suburb,
		CityDaneCode:  dto.CityDaneCode,
		PostalCode:    dto.PostalCode,
		Street:        dto.Street,
		Latitude:      dto.Latitude,
		Longitude:     dto.Longitude,
	}

	return uc.repo.Create(ctx, warehouse)
}
