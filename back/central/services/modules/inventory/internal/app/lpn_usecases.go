package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) CreateLPN(ctx context.Context, dto request.CreateLPNDTO) (*entities.LicensePlate, error) {
	dup, err := uc.repo.LPNExistsByCode(ctx, dto.BusinessID, dto.Code, nil)
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, domainerrors.ErrDuplicateLPNCode
	}
	lpnType := dto.LpnType
	if lpnType == "" {
		lpnType = "pallet"
	}
	lpn := &entities.LicensePlate{
		BusinessID:        dto.BusinessID,
		Code:              dto.Code,
		LpnType:           lpnType,
		CurrentLocationID: dto.LocationID,
		Status:            "active",
	}
	return uc.repo.CreateLPN(ctx, lpn)
}

func (uc *useCase) GetLPN(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
	return uc.repo.GetLPNByID(ctx, businessID, id)
}

func (uc *useCase) ListLPNs(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListLPNs(ctx, params)
}

func (uc *useCase) UpdateLPN(ctx context.Context, dto request.UpdateLPNDTO) (*entities.LicensePlate, error) {
	existing, err := uc.repo.GetLPNByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}
	if dto.Code != "" && dto.Code != existing.Code {
		dup, err := uc.repo.LPNExistsByCode(ctx, dto.BusinessID, dto.Code, &existing.ID)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, domainerrors.ErrDuplicateLPNCode
		}
		existing.Code = dto.Code
	}
	if dto.LpnType != "" {
		existing.LpnType = dto.LpnType
	}
	if dto.LocationID != nil {
		existing.CurrentLocationID = dto.LocationID
	}
	if dto.Status != "" {
		existing.Status = dto.Status
	}
	return uc.repo.UpdateLPN(ctx, existing)
}

func (uc *useCase) DeleteLPN(ctx context.Context, businessID, id uint) error {
	return uc.repo.DeleteLPN(ctx, businessID, id)
}

func (uc *useCase) AddToLPN(ctx context.Context, dto request.AddToLPNDTO) (*entities.LicensePlateLine, error) {
	if dto.Qty <= 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}
	lpn, err := uc.repo.GetLPNByID(ctx, dto.BusinessID, dto.LpnID)
	if err != nil {
		return nil, err
	}
	if lpn.Status == "dissolved" {
		return nil, domainerrors.ErrLPNDissolved
	}
	line := &entities.LicensePlateLine{
		LpnID:      dto.LpnID,
		BusinessID: dto.BusinessID,
		ProductID:  dto.ProductID,
		LotID:      dto.LotID,
		SerialID:   dto.SerialID,
		Qty:        dto.Qty,
	}
	return uc.repo.AddLPNLine(ctx, line)
}

func (uc *useCase) MoveLPN(ctx context.Context, dto request.MoveLPNDTO) (*entities.LicensePlate, error) {
	existing, err := uc.repo.GetLPNByID(ctx, dto.BusinessID, dto.LpnID)
	if err != nil {
		return nil, err
	}
	existing.CurrentLocationID = &dto.NewLocationID
	return uc.repo.UpdateLPN(ctx, existing)
}

func (uc *useCase) DissolveLPN(ctx context.Context, dto request.DissolveLPNDTO) error {
	existing, err := uc.repo.GetLPNByID(ctx, dto.BusinessID, dto.LpnID)
	if err != nil {
		return err
	}
	if existing.Status == "dissolved" {
		return domainerrors.ErrLPNDissolved
	}
	return uc.repo.DissolveLPN(ctx, dto.BusinessID, dto.LpnID)
}

func (uc *useCase) MergeLPN(ctx context.Context, dto request.MergeLPNDTO) (*entities.LicensePlate, error) {
	source, err := uc.repo.GetLPNByID(ctx, dto.BusinessID, dto.SourceLpnID)
	if err != nil {
		return nil, err
	}
	target, err := uc.repo.GetLPNByID(ctx, dto.BusinessID, dto.TargetLpnID)
	if err != nil {
		return nil, err
	}
	lines, err := uc.repo.ListLPNLines(ctx, source.ID)
	if err != nil {
		return nil, err
	}
	for _, line := range lines {
		line.ID = 0
		line.LpnID = target.ID
		line.BusinessID = dto.BusinessID
		_, _ = uc.repo.AddLPNLine(ctx, &line)
	}
	_ = uc.repo.DissolveLPN(ctx, dto.BusinessID, source.ID)
	return uc.repo.GetLPNByID(ctx, dto.BusinessID, target.ID)
}
