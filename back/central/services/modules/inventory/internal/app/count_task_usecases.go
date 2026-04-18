package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) GenerateCountTask(ctx context.Context, dto request.GenerateCountTaskDTO) (*response.GenerateCountTaskResult, error) {
	plan, err := uc.repo.GetCountPlanByID(ctx, dto.BusinessID, dto.PlanID)
	if err != nil {
		return nil, err
	}
	scopeType := dto.ScopeType
	if scopeType == "" {
		scopeType = plan.Strategy
	}
	task := &entities.CycleCountTask{
		PlanID:      plan.ID,
		BusinessID:  plan.BusinessID,
		WarehouseID: plan.WarehouseID,
		ScopeType:   scopeType,
		ScopeID:     dto.ScopeID,
		Status:      "pending",
	}
	created, err := uc.repo.CreateCountTask(ctx, task)
	if err != nil {
		return nil, err
	}
	lines, err := uc.repo.GenerateCountLinesForTask(ctx, created, plan.Strategy)
	if err != nil {
		return nil, err
	}
	return &response.GenerateCountTaskResult{Task: *created, Lines: lines}, nil
}

func (uc *useCase) ListCountTasks(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListCountTasks(ctx, params)
}

func (uc *useCase) GetCountTask(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
	return uc.repo.GetCountTaskByID(ctx, businessID, id)
}

func (uc *useCase) StartCountTask(ctx context.Context, dto request.StartCountTaskDTO) (*entities.CycleCountTask, error) {
	existing, err := uc.repo.GetCountTaskByID(ctx, dto.BusinessID, dto.TaskID)
	if err != nil {
		return nil, err
	}
	if existing.Status == "completed" || existing.Status == "cancelled" {
		return nil, domainerrors.ErrCountTaskClosed
	}
	now := time.Now()
	existing.Status = "in_progress"
	existing.AssignedToID = &dto.UserID
	existing.StartedAt = &now
	return uc.repo.UpdateCountTask(ctx, existing)
}

func (uc *useCase) FinishCountTask(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
	existing, err := uc.repo.GetCountTaskByID(ctx, businessID, id)
	if err != nil {
		return nil, err
	}
	if existing.Status == "completed" {
		return nil, domainerrors.ErrCountTaskClosed
	}
	now := time.Now()
	existing.Status = "completed"
	existing.FinishedAt = &now
	return uc.repo.UpdateCountTask(ctx, existing)
}

func (uc *useCase) ListCountLines(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListCountLines(ctx, params)
}

func (uc *useCase) SubmitCountLine(ctx context.Context, dto request.SubmitCountLineDTO) (*response.SubmitCountLineResult, error) {
	line, err := uc.repo.GetCountLineByID(ctx, dto.BusinessID, dto.LineID)
	if err != nil {
		return nil, err
	}
	if line.Status == "submitted" || line.Status == "resolved" {
		return nil, domainerrors.ErrCountLineSubmitted
	}
	counted := dto.CountedQty
	line.CountedQty = &counted
	line.Variance = counted - line.ExpectedQty
	if line.Variance == 0 {
		line.Status = "resolved"
	} else {
		line.Status = "submitted"
	}
	updated, err := uc.repo.UpdateCountLine(ctx, line)
	if err != nil {
		return nil, err
	}
	result := &response.SubmitCountLineResult{Line: *updated}
	if updated.Variance != 0 {
		disc := &entities.InventoryDiscrepancy{
			TaskID:     updated.TaskID,
			LineID:     updated.ID,
			BusinessID: updated.BusinessID,
			Status:     "open",
		}
		createdDisc, err := uc.repo.CreateDiscrepancy(ctx, disc)
		if err == nil {
			result.Discrepancy = createdDisc
		}
	}
	return result, nil
}
