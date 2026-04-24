package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (m *UseCaseMock) CreateCountPlan(ctx context.Context, dto request.CreateCountPlanDTO) (*entities.CycleCountPlan, error) {
	if m.CreateCountPlanFn != nil {
		return m.CreateCountPlanFn(ctx, dto)
	}
	return &entities.CycleCountPlan{}, nil
}

func (m *UseCaseMock) GetCountPlan(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
	if m.GetCountPlanFn != nil {
		return m.GetCountPlanFn(ctx, businessID, id)
	}
	return &entities.CycleCountPlan{ID: id, BusinessID: businessID}, nil
}

func (m *UseCaseMock) ListCountPlans(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error) {
	if m.ListCountPlansFn != nil {
		return m.ListCountPlansFn(ctx, params)
	}
	return []entities.CycleCountPlan{}, 0, nil
}

func (m *UseCaseMock) UpdateCountPlan(ctx context.Context, dto request.UpdateCountPlanDTO) (*entities.CycleCountPlan, error) {
	if m.UpdateCountPlanFn != nil {
		return m.UpdateCountPlanFn(ctx, dto)
	}
	return &entities.CycleCountPlan{}, nil
}

func (m *UseCaseMock) DeleteCountPlan(ctx context.Context, businessID, id uint) error {
	if m.DeleteCountPlanFn != nil {
		return m.DeleteCountPlanFn(ctx, businessID, id)
	}
	return nil
}

func (m *UseCaseMock) GenerateCountTask(ctx context.Context, dto request.GenerateCountTaskDTO) (*response.GenerateCountTaskResult, error) {
	if m.GenerateCountTaskFn != nil {
		return m.GenerateCountTaskFn(ctx, dto)
	}
	return &response.GenerateCountTaskResult{}, nil
}

func (m *UseCaseMock) ListCountTasks(ctx context.Context, params dtos.ListCycleCountTasksParams) ([]entities.CycleCountTask, int64, error) {
	if m.ListCountTasksFn != nil {
		return m.ListCountTasksFn(ctx, params)
	}
	return []entities.CycleCountTask{}, 0, nil
}

func (m *UseCaseMock) GetCountTask(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
	if m.GetCountTaskFn != nil {
		return m.GetCountTaskFn(ctx, businessID, id)
	}
	return &entities.CycleCountTask{ID: id, BusinessID: businessID}, nil
}

func (m *UseCaseMock) StartCountTask(ctx context.Context, dto request.StartCountTaskDTO) (*entities.CycleCountTask, error) {
	if m.StartCountTaskFn != nil {
		return m.StartCountTaskFn(ctx, dto)
	}
	return &entities.CycleCountTask{}, nil
}

func (m *UseCaseMock) FinishCountTask(ctx context.Context, businessID, id uint) (*entities.CycleCountTask, error) {
	if m.FinishCountTaskFn != nil {
		return m.FinishCountTaskFn(ctx, businessID, id)
	}
	return &entities.CycleCountTask{}, nil
}

func (m *UseCaseMock) ListCountLines(ctx context.Context, params dtos.ListCycleCountLinesParams) ([]entities.CycleCountLine, int64, error) {
	if m.ListCountLinesFn != nil {
		return m.ListCountLinesFn(ctx, params)
	}
	return []entities.CycleCountLine{}, 0, nil
}

func (m *UseCaseMock) SubmitCountLine(ctx context.Context, dto request.SubmitCountLineDTO) (*response.SubmitCountLineResult, error) {
	if m.SubmitCountLineFn != nil {
		return m.SubmitCountLineFn(ctx, dto)
	}
	return &response.SubmitCountLineResult{}, nil
}

func (m *UseCaseMock) ListDiscrepancies(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error) {
	if m.ListDiscrepanciesFn != nil {
		return m.ListDiscrepanciesFn(ctx, params)
	}
	return []entities.InventoryDiscrepancy{}, 0, nil
}

func (m *UseCaseMock) GetDiscrepancy(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
	if m.GetDiscrepancyFn != nil {
		return m.GetDiscrepancyFn(ctx, businessID, id)
	}
	return &entities.InventoryDiscrepancy{ID: id, BusinessID: businessID}, nil
}

func (m *UseCaseMock) ApproveDiscrepancy(ctx context.Context, dto request.ApproveDiscrepancyDTO) (*entities.InventoryDiscrepancy, error) {
	if m.ApproveDiscrepancyFn != nil {
		return m.ApproveDiscrepancyFn(ctx, dto)
	}
	return &entities.InventoryDiscrepancy{}, nil
}

func (m *UseCaseMock) RejectDiscrepancy(ctx context.Context, dto request.RejectDiscrepancyDTO) (*entities.InventoryDiscrepancy, error) {
	if m.RejectDiscrepancyFn != nil {
		return m.RejectDiscrepancyFn(ctx, dto)
	}
	return &entities.InventoryDiscrepancy{}, nil
}

func (m *UseCaseMock) ExportKardex(ctx context.Context, dto request.KardexExportDTO) (*response.KardexExportResult, error) {
	if m.ExportKardexFn != nil {
		return m.ExportKardexFn(ctx, dto)
	}
	return &response.KardexExportResult{}, nil
}
