package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) CreateSerial(ctx context.Context, dto request.CreateSerialDTO) (*entities.InventorySerial, error) {
	if _, _, _, err := uc.repo.GetProductByID(ctx, dto.ProductID, dto.BusinessID); err != nil {
		return nil, domainerrors.ErrProductNotFound
	}

	dup, err := uc.repo.SerialExists(ctx, dto.BusinessID, dto.ProductID, dto.SerialNumber, nil)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, domainerrors.ErrDuplicateSerial
	}

	var stateID *uint
	if dto.StateCode != "" {
		state, err := uc.repo.GetInventoryStateByCode(ctx, dto.StateCode)
		if err != nil {
			return nil, err
		}
		stateID = &state.ID
	}

	now := time.Now()
	serial := &entities.InventorySerial{
		BusinessID:        dto.BusinessID,
		ProductID:         dto.ProductID,
		SerialNumber:      dto.SerialNumber,
		LotID:             dto.LotID,
		CurrentLocationID: dto.LocationID,
		CurrentStateID:    stateID,
		ReceivedAt:        &now,
	}
	return uc.repo.CreateSerial(ctx, serial)
}

func (uc *useCase) GetSerial(ctx context.Context, businessID, serialID uint) (*entities.InventorySerial, error) {
	return uc.repo.GetSerialByID(ctx, businessID, serialID)
}

func (uc *useCase) ListSerials(ctx context.Context, params dtos.ListSerialsParams) ([]entities.InventorySerial, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListSerials(ctx, params)
}

func (uc *useCase) UpdateSerial(ctx context.Context, dto request.UpdateSerialDTO) (*entities.InventorySerial, error) {
	existing, err := uc.repo.GetSerialByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.LotID != nil {
		existing.LotID = dto.LotID
	}
	if dto.LocationID != nil {
		existing.CurrentLocationID = dto.LocationID
	}
	if dto.StateCode != "" {
		state, err := uc.repo.GetInventoryStateByCode(ctx, dto.StateCode)
		if err != nil {
			return nil, err
		}
		existing.CurrentStateID = &state.ID
	}

	return uc.repo.UpdateSerial(ctx, existing)
}
