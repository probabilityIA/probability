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

func (uc *useCase) CreateReplenishmentTask(ctx context.Context, dto request.CreateReplenishmentTaskDTO) (*entities.ReplenishmentTask, error) {
	if dto.Quantity <= 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}
	triggered := dto.TriggeredBy
	if triggered == "" {
		triggered = "manual"
	}
	task := &entities.ReplenishmentTask{
		BusinessID:     dto.BusinessID,
		ProductID:      dto.ProductID,
		WarehouseID:    dto.WarehouseID,
		FromLocationID: dto.FromLocationID,
		ToLocationID:   dto.ToLocationID,
		Quantity:       dto.Quantity,
		Status:         "pending",
		TriggeredBy:    triggered,
		Notes:          dto.Notes,
	}
	return uc.repo.CreateReplenishmentTask(ctx, task)
}

func (uc *useCase) ListReplenishmentTasks(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListReplenishmentTasks(ctx, params)
}

func (uc *useCase) AssignReplenishment(ctx context.Context, dto request.AssignReplenishmentDTO) (*entities.ReplenishmentTask, error) {
	existing, err := uc.repo.GetReplenishmentTaskByID(ctx, dto.BusinessID, dto.TaskID)
	if err != nil {
		return nil, err
	}
	if existing.Status == "completed" || existing.Status == "cancelled" {
		return nil, domainerrors.ErrReplenishmentClosed
	}
	now := time.Now()
	existing.AssignedToID = &dto.UserID
	existing.AssignedAt = &now
	existing.Status = "in_progress"
	return uc.repo.UpdateReplenishmentTask(ctx, existing)
}

func (uc *useCase) CompleteReplenishment(ctx context.Context, dto request.CompleteReplenishmentDTO) (*entities.ReplenishmentTask, error) {
	existing, err := uc.repo.GetReplenishmentTaskByID(ctx, dto.BusinessID, dto.TaskID)
	if err != nil {
		return nil, err
	}
	if existing.Status == "completed" || existing.Status == "cancelled" {
		return nil, domainerrors.ErrReplenishmentClosed
	}
	now := time.Now()
	existing.Status = "completed"
	existing.CompletedAt = &now
	if dto.Notes != "" {
		existing.Notes = dto.Notes
	}
	return uc.repo.UpdateReplenishmentTask(ctx, existing)
}

func (uc *useCase) CancelReplenishment(ctx context.Context, businessID, taskID uint, reason string) (*entities.ReplenishmentTask, error) {
	existing, err := uc.repo.GetReplenishmentTaskByID(ctx, businessID, taskID)
	if err != nil {
		return nil, err
	}
	if existing.Status == "completed" {
		return nil, domainerrors.ErrReplenishmentClosed
	}
	existing.Status = "cancelled"
	if reason != "" {
		existing.Notes = reason
	}
	return uc.repo.UpdateReplenishmentTask(ctx, existing)
}

func (uc *useCase) DetectReplenishmentNeeds(ctx context.Context, businessID uint) (*response.ReplenishmentDetectResult, error) {
	candidates, err := uc.repo.DetectReplenishmentCandidates(ctx, businessID)
	if err != nil {
		return nil, err
	}
	result := &response.ReplenishmentDetectResult{}
	for i := range candidates {
		created, err := uc.repo.CreateReplenishmentTask(ctx, &candidates[i])
		if err != nil {
			continue
		}
		result.Tasks = append(result.Tasks, *created)
		result.Created++
	}
	return result, nil
}
