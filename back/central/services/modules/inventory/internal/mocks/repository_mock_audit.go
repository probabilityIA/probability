package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (m *RepositoryMock) CreateCountPlan(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
	if m.CreateCountPlanFn != nil {
		return m.CreateCountPlanFn(ctx, p)
	}
	return p, nil
}

func (m *RepositoryMock) GetCountPlanByID(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
	if m.GetCountPlanByIDFn != nil {
		return m.GetCountPlanByIDFn(ctx, businessID, id)
	}
	return &entities.CycleCountPlan{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListCountPlans(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error) {
	if m.ListCountPlansFn != nil {
		return m.ListCountPlansFn(ctx, params)
	}
	return []entities.CycleCountPlan{}, 0, nil
}

func (m *RepositoryMock) UpdateCountPlan(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
	if m.UpdateCountPlanFn != nil {
		return m.UpdateCountPlanFn(ctx, p)
	}
	return p, nil
}

func (m *RepositoryMock) DeleteCountPlan(ctx context.Context, businessID, id uint) error {
	if m.DeleteCountPlanFn != nil {
		return m.DeleteCountPlanFn(ctx, businessID, id)
	}
	return nil
}

func (m *RepositoryMock) CreateCountTask(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
	if m.CreateCountTaskFn != nil {
		return m.CreateCountTaskFn(ctx, t)
	}
	return t, nil
}

func (m *RepositoryMock) GetCountTaskByID(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
	if m.GetCountTaskByIDFn != nil {
		return m.GetCountTaskByIDFn(ctx, businessID, id)
	}
	return &entities.CycleCountTask{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListCountTasks(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error) {
	if m.ListCountTasksFn != nil {
		return m.ListCountTasksFn(ctx, params)
	}
	return []entities.CycleCountTask{}, 0, nil
}

func (m *RepositoryMock) UpdateCountTask(ctx context.Context, t *entities.CycleCountTask) (*entities.CycleCountTask, error) {
	if m.UpdateCountTaskFn != nil {
		return m.UpdateCountTaskFn(ctx, t)
	}
	return t, nil
}

func (m *RepositoryMock) GenerateCountLinesForTask(ctx context.Context, task *entities.CycleCountTask, strategy string) ([]entities.CycleCountLine, error) {
	if m.GenerateCountLinesForTaskFn != nil {
		return m.GenerateCountLinesForTaskFn(ctx, task, strategy)
	}
	return []entities.CycleCountLine{}, nil
}

func (m *RepositoryMock) CreateCountLine(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
	if m.CreateCountLineFn != nil {
		return m.CreateCountLineFn(ctx, line)
	}
	return line, nil
}

func (m *RepositoryMock) GetCountLineByID(ctx context.Context, businessID, id uint) (*entities.CycleCountLine, error) {
	if m.GetCountLineByIDFn != nil {
		return m.GetCountLineByIDFn(ctx, businessID, id)
	}
	return &entities.CycleCountLine{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListCountLines(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error) {
	if m.ListCountLinesFn != nil {
		return m.ListCountLinesFn(ctx, params)
	}
	return []entities.CycleCountLine{}, 0, nil
}

func (m *RepositoryMock) UpdateCountLine(ctx context.Context, line *entities.CycleCountLine) (*entities.CycleCountLine, error) {
	if m.UpdateCountLineFn != nil {
		return m.UpdateCountLineFn(ctx, line)
	}
	return line, nil
}

func (m *RepositoryMock) CreateDiscrepancy(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
	if m.CreateDiscrepancyFn != nil {
		return m.CreateDiscrepancyFn(ctx, d)
	}
	return d, nil
}

func (m *RepositoryMock) GetDiscrepancyByID(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
	if m.GetDiscrepancyByIDFn != nil {
		return m.GetDiscrepancyByIDFn(ctx, businessID, id)
	}
	return &entities.InventoryDiscrepancy{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListDiscrepancies(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error) {
	if m.ListDiscrepanciesFn != nil {
		return m.ListDiscrepanciesFn(ctx, params)
	}
	return []entities.InventoryDiscrepancy{}, 0, nil
}

func (m *RepositoryMock) UpdateDiscrepancy(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
	if m.UpdateDiscrepancyFn != nil {
		return m.UpdateDiscrepancyFn(ctx, d)
	}
	return d, nil
}

func (m *RepositoryMock) ApproveDiscrepancyTx(ctx context.Context, params dtos.ApproveDiscrepancyTxParams) (*entities.InventoryDiscrepancy, error) {
	if m.ApproveDiscrepancyTxFn != nil {
		return m.ApproveDiscrepancyTxFn(ctx, params)
	}
	return &entities.InventoryDiscrepancy{}, nil
}

func (m *RepositoryMock) GetKardex(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error) {
	if m.GetKardexFn != nil {
		return m.GetKardexFn(ctx, params)
	}
	return []entities.KardexEntry{}, nil
}
