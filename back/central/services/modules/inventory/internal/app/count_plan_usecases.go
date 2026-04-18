package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (uc *useCase) CreateCountPlan(ctx context.Context, dto request.CreateCountPlanDTO) (*entities.CycleCountPlan, error) {
	strategy := dto.Strategy
	if strategy == "" {
		strategy = "abc"
	}
	if dto.FrequencyDays <= 0 {
		dto.FrequencyDays = 30
	}
	plan := &entities.CycleCountPlan{
		BusinessID:    dto.BusinessID,
		WarehouseID:   dto.WarehouseID,
		Name:          dto.Name,
		Strategy:      strategy,
		FrequencyDays: dto.FrequencyDays,
		NextRunAt:     dto.NextRunAt,
		IsActive:      dto.IsActive,
	}
	return uc.repo.CreateCountPlan(ctx, plan)
}

func (uc *useCase) GetCountPlan(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
	return uc.repo.GetCountPlanByID(ctx, businessID, id)
}

func (uc *useCase) ListCountPlans(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListCountPlans(ctx, params)
}

func (uc *useCase) UpdateCountPlan(ctx context.Context, dto request.UpdateCountPlanDTO) (*entities.CycleCountPlan, error) {
	existing, err := uc.repo.GetCountPlanByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}
	if dto.WarehouseID != nil {
		existing.WarehouseID = *dto.WarehouseID
	}
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.Strategy != "" {
		existing.Strategy = dto.Strategy
	}
	if dto.FrequencyDays != nil {
		existing.FrequencyDays = *dto.FrequencyDays
	}
	if dto.NextRunAt != nil {
		existing.NextRunAt = dto.NextRunAt
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	return uc.repo.UpdateCountPlan(ctx, existing)
}

func (uc *useCase) DeleteCountPlan(ctx context.Context, businessID, id uint) error {
	return uc.repo.DeleteCountPlan(ctx, businessID, id)
}
