package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) CreateLot(ctx context.Context, dto request.CreateLotDTO) (*entities.InventoryLot, error) {
	if _, _, _, err := uc.repo.GetProductByID(ctx, dto.ProductID, dto.BusinessID); err != nil {
		return nil, domainerrors.ErrProductNotFound
	}

	dup, err := uc.repo.LotExistsByCode(ctx, dto.BusinessID, dto.ProductID, dto.LotCode, nil)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, domainerrors.ErrDuplicateLotCode
	}

	status := dto.Status
	if status == "" {
		status = "active"
	}

	lot := &entities.InventoryLot{
		BusinessID:      dto.BusinessID,
		ProductID:       dto.ProductID,
		LotCode:         dto.LotCode,
		ManufactureDate: dto.ManufactureDate,
		ExpirationDate:  dto.ExpirationDate,
		ReceivedAt:      dto.ReceivedAt,
		SupplierID:      dto.SupplierID,
		Status:          status,
	}
	return uc.repo.CreateLot(ctx, lot)
}

func (uc *useCase) GetLot(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
	return uc.repo.GetLotByID(ctx, businessID, lotID)
}

func (uc *useCase) ListLots(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListLots(ctx, params)
}

func (uc *useCase) UpdateLot(ctx context.Context, dto request.UpdateLotDTO) (*entities.InventoryLot, error) {
	existing, err := uc.repo.GetLotByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}

	if dto.LotCode != "" && dto.LotCode != existing.LotCode {
		dup, err := uc.repo.LotExistsByCode(ctx, dto.BusinessID, existing.ProductID, dto.LotCode, &existing.ID)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, domainerrors.ErrDuplicateLotCode
		}
		existing.LotCode = dto.LotCode
	}
	if dto.ManufactureDate != nil {
		existing.ManufactureDate = dto.ManufactureDate
	}
	if dto.ExpirationDate != nil {
		existing.ExpirationDate = dto.ExpirationDate
	}
	if dto.ReceivedAt != nil {
		existing.ReceivedAt = dto.ReceivedAt
	}
	if dto.SupplierID != nil {
		existing.SupplierID = dto.SupplierID
	}
	if dto.Status != "" {
		existing.Status = dto.Status
	}

	return uc.repo.UpdateLot(ctx, existing)
}

func (uc *useCase) DeleteLot(ctx context.Context, businessID, lotID uint) error {
	return uc.repo.DeleteLot(ctx, businessID, lotID)
}
