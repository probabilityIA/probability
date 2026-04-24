package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
)

func (u *UseCase) CreateRack(ctx context.Context, dto request.CreateRackDTO) (*entities.WarehouseRack, error) {
	if _, err := u.repo.GetAisleByID(ctx, dto.BusinessID, dto.AisleID); err != nil {
		return nil, err
	}

	dup, err := u.repo.RackExistsByCode(ctx, dto.AisleID, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, domainerrors.ErrDuplicateRackCode
	}

	rack := &entities.WarehouseRack{
		AisleID:     dto.AisleID,
		BusinessID:  dto.BusinessID,
		Code:        dto.Code,
		Name:        dto.Name,
		LevelsCount: dto.LevelsCount,
		IsActive:    dto.IsActive,
	}
	return u.repo.CreateRack(ctx, rack)
}

func (u *UseCase) GetRack(ctx context.Context, businessID, rackID uint) (*entities.WarehouseRack, error) {
	return u.repo.GetRackByID(ctx, businessID, rackID)
}

func (u *UseCase) ListRacks(ctx context.Context, params dtos.ListRacksParams) ([]entities.WarehouseRack, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return u.repo.ListRacks(ctx, params)
}

func (u *UseCase) UpdateRack(ctx context.Context, dto request.UpdateRackDTO) (*entities.WarehouseRack, error) {
	existing, err := u.repo.GetRackByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.Code != "" && dto.Code != existing.Code {
		dup, err := u.repo.RackExistsByCode(ctx, existing.AisleID, dto.Code, &existing.ID)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, domainerrors.ErrDuplicateRackCode
		}
		existing.Code = dto.Code
	}
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.LevelsCount != nil {
		existing.LevelsCount = *dto.LevelsCount
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}

	return u.repo.UpdateRack(ctx, existing)
}

func (u *UseCase) DeleteRack(ctx context.Context, businessID, rackID uint) error {
	return u.repo.DeleteRack(ctx, businessID, rackID)
}
