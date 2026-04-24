package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
)

func (u *UseCase) CreateRackLevel(ctx context.Context, dto request.CreateRackLevelDTO) (*entities.WarehouseRackLevel, error) {
	if _, err := u.repo.GetRackByID(ctx, dto.BusinessID, dto.RackID); err != nil {
		return nil, err
	}

	dup, err := u.repo.RackLevelExistsByCode(ctx, dto.RackID, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, domainerrors.ErrDuplicateLevelCode
	}

	level := &entities.WarehouseRackLevel{
		RackID:     dto.RackID,
		BusinessID: dto.BusinessID,
		Code:       dto.Code,
		Ordinal:    dto.Ordinal,
		IsActive:   dto.IsActive,
	}
	return u.repo.CreateRackLevel(ctx, level)
}

func (u *UseCase) GetRackLevel(ctx context.Context, businessID, levelID uint) (*entities.WarehouseRackLevel, error) {
	return u.repo.GetRackLevelByID(ctx, businessID, levelID)
}

func (u *UseCase) ListRackLevels(ctx context.Context, params dtos.ListRackLevelsParams) ([]entities.WarehouseRackLevel, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return u.repo.ListRackLevels(ctx, params)
}

func (u *UseCase) UpdateRackLevel(ctx context.Context, dto request.UpdateRackLevelDTO) (*entities.WarehouseRackLevel, error) {
	existing, err := u.repo.GetRackLevelByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.Code != "" && dto.Code != existing.Code {
		dup, err := u.repo.RackLevelExistsByCode(ctx, existing.RackID, dto.Code, &existing.ID)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, domainerrors.ErrDuplicateLevelCode
		}
		existing.Code = dto.Code
	}
	if dto.Ordinal != nil {
		existing.Ordinal = *dto.Ordinal
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}

	return u.repo.UpdateRackLevel(ctx, existing)
}

func (u *UseCase) DeleteRackLevel(ctx context.Context, businessID, levelID uint) error {
	return u.repo.DeleteRackLevel(ctx, businessID, levelID)
}
