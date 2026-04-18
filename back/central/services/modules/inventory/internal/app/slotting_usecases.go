package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (uc *useCase) RunSlotting(ctx context.Context, dto request.RunSlottingDTO) (*response.SlottingRunResult, error) {
	period := dto.Period
	if period == "" {
		period = "30d"
	}
	if err := uc.repo.ComputeVelocities(ctx, dto.BusinessID, dto.WarehouseID, period); err != nil {
		return nil, err
	}
	velocities, err := uc.repo.ListVelocities(ctx, dtos.ListVelocityParams{
		BusinessID:  dto.BusinessID,
		WarehouseID: dto.WarehouseID,
		Period:      period,
		Limit:       100,
	})
	if err != nil {
		return nil, err
	}
	return &response.SlottingRunResult{
		BusinessID:   dto.BusinessID,
		WarehouseID:  dto.WarehouseID,
		Period:       period,
		TotalScanned: len(velocities),
		Velocities:   velocities,
	}, nil
}

func (uc *useCase) ListVelocities(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error) {
	return uc.repo.ListVelocities(ctx, params)
}
