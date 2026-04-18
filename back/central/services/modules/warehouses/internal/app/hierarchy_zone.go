package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/errors"
)

func (u *UseCase) CreateZone(ctx context.Context, dto request.CreateZoneDTO) (*entities.WarehouseZone, error) {
	exists, err := u.repo.WarehouseExists(ctx, dto.BusinessID, dto.WarehouseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrWarehouseNotFound
	}

	dup, err := u.repo.ZoneExistsByCode(ctx, dto.WarehouseID, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, domainerrors.ErrDuplicateZoneCode
	}

	zone := &entities.WarehouseZone{
		WarehouseID: dto.WarehouseID,
		BusinessID:  dto.BusinessID,
		Code:        dto.Code,
		Name:        dto.Name,
		Purpose:     dto.Purpose,
		ColorHex:    dto.ColorHex,
		IsActive:    dto.IsActive,
	}
	return u.repo.CreateZone(ctx, zone)
}

func (u *UseCase) GetZone(ctx context.Context, businessID, zoneID uint) (*entities.WarehouseZone, error) {
	return u.repo.GetZoneByID(ctx, businessID, zoneID)
}

func (u *UseCase) ListZones(ctx context.Context, params dtos.ListZonesParams) ([]entities.WarehouseZone, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return u.repo.ListZones(ctx, params)
}

func (u *UseCase) UpdateZone(ctx context.Context, dto request.UpdateZoneDTO) (*entities.WarehouseZone, error) {
	existing, err := u.repo.GetZoneByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.Code != "" && dto.Code != existing.Code {
		dup, err := u.repo.ZoneExistsByCode(ctx, existing.WarehouseID, dto.Code, &existing.ID)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, domainerrors.ErrDuplicateZoneCode
		}
		existing.Code = dto.Code
	}
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.Purpose != "" {
		existing.Purpose = dto.Purpose
	}
	if dto.ColorHex != "" {
		existing.ColorHex = dto.ColorHex
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}

	return u.repo.UpdateZone(ctx, existing)
}

func (u *UseCase) DeleteZone(ctx context.Context, businessID, zoneID uint) error {
	return u.repo.DeleteZone(ctx, businessID, zoneID)
}
