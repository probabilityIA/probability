package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
)

func (u *UseCase) CreateAisle(ctx context.Context, dto request.CreateAisleDTO) (*entities.WarehouseAisle, error) {
	zone, err := u.repo.GetZoneByID(ctx, dto.BusinessID, dto.ZoneID)
	if err != nil {
		return nil, err
	}
	_ = zone

	dup, err := u.repo.AisleExistsByCode(ctx, dto.ZoneID, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, domainerrors.ErrDuplicateAisleCode
	}

	aisle := &entities.WarehouseAisle{
		ZoneID:     dto.ZoneID,
		BusinessID: dto.BusinessID,
		Code:       dto.Code,
		Name:       dto.Name,
		IsActive:   dto.IsActive,
	}
	return u.repo.CreateAisle(ctx, aisle)
}

func (u *UseCase) GetAisle(ctx context.Context, businessID, aisleID uint) (*entities.WarehouseAisle, error) {
	return u.repo.GetAisleByID(ctx, businessID, aisleID)
}

func (u *UseCase) ListAisles(ctx context.Context, params dtos.ListAislesParams) ([]entities.WarehouseAisle, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return u.repo.ListAisles(ctx, params)
}

func (u *UseCase) UpdateAisle(ctx context.Context, dto request.UpdateAisleDTO) (*entities.WarehouseAisle, error) {
	existing, err := u.repo.GetAisleByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.Code != "" && dto.Code != existing.Code {
		dup, err := u.repo.AisleExistsByCode(ctx, existing.ZoneID, dto.Code, &existing.ID)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, domainerrors.ErrDuplicateAisleCode
		}
		existing.Code = dto.Code
	}
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}

	return u.repo.UpdateAisle(ctx, existing)
}

func (u *UseCase) DeleteAisle(ctx context.Context, businessID, aisleID uint) error {
	return u.repo.DeleteAisle(ctx, businessID, aisleID)
}
