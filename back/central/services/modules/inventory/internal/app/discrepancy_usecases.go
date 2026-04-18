package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) ListDiscrepancies(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListDiscrepancies(ctx, params)
}

func (uc *useCase) GetDiscrepancy(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
	return uc.repo.GetDiscrepancyByID(ctx, businessID, id)
}

func (uc *useCase) ApproveDiscrepancy(ctx context.Context, dto request.ApproveDiscrepancyDTO) (*entities.InventoryDiscrepancy, error) {
	movTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, "count_adjustment")
	if err != nil {
		return nil, err
	}
	return uc.repo.ApproveDiscrepancyTx(ctx, dtos.ApproveDiscrepancyTxParams{
		BusinessID:     dto.BusinessID,
		DiscrepancyID:  dto.DiscrepancyID,
		ReviewerID:     dto.ReviewerID,
		Notes:          dto.Notes,
		MovementTypeID: movTypeID,
	})
}

func (uc *useCase) RejectDiscrepancy(ctx context.Context, dto request.RejectDiscrepancyDTO) (*entities.InventoryDiscrepancy, error) {
	existing, err := uc.repo.GetDiscrepancyByID(ctx, dto.BusinessID, dto.DiscrepancyID)
	if err != nil {
		return nil, err
	}
	if existing.Status == "approved" || existing.Status == "rejected" {
		return nil, domainerrors.ErrDiscrepancyResolved
	}
	now := time.Now()
	existing.Status = "rejected"
	existing.ReviewedByID = &dto.ReviewerID
	existing.ReviewedAt = &now
	if dto.Reason != "" {
		existing.Notes = dto.Reason
	}
	return uc.repo.UpdateDiscrepancy(ctx, existing)
}
