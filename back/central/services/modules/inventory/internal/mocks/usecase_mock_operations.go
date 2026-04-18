package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (m *UseCaseMock) CreatePutawayRule(ctx context.Context, dto request.CreatePutawayRuleDTO) (*entities.PutawayRule, error) {
	if m.CreatePutawayRuleFn != nil {
		return m.CreatePutawayRuleFn(ctx, dto)
	}
	return &entities.PutawayRule{}, nil
}

func (m *UseCaseMock) ListPutawayRules(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error) {
	if m.ListPutawayRulesFn != nil {
		return m.ListPutawayRulesFn(ctx, params)
	}
	return []entities.PutawayRule{}, 0, nil
}

func (m *UseCaseMock) UpdatePutawayRule(ctx context.Context, dto request.UpdatePutawayRuleDTO) (*entities.PutawayRule, error) {
	if m.UpdatePutawayRuleFn != nil {
		return m.UpdatePutawayRuleFn(ctx, dto)
	}
	return &entities.PutawayRule{}, nil
}

func (m *UseCaseMock) DeletePutawayRule(ctx context.Context, businessID, ruleID uint) error {
	if m.DeletePutawayRuleFn != nil {
		return m.DeletePutawayRuleFn(ctx, businessID, ruleID)
	}
	return nil
}

func (m *UseCaseMock) SuggestPutaway(ctx context.Context, dto request.PutawaySuggestDTO) (*response.PutawaySuggestResult, error) {
	if m.SuggestPutawayFn != nil {
		return m.SuggestPutawayFn(ctx, dto)
	}
	return &response.PutawaySuggestResult{}, nil
}

func (m *UseCaseMock) ConfirmPutaway(ctx context.Context, dto request.ConfirmPutawayDTO) (*entities.PutawaySuggestion, error) {
	if m.ConfirmPutawayFn != nil {
		return m.ConfirmPutawayFn(ctx, dto)
	}
	return &entities.PutawaySuggestion{}, nil
}

func (m *UseCaseMock) ListPutawaySuggestions(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error) {
	if m.ListPutawaySuggestionsFn != nil {
		return m.ListPutawaySuggestionsFn(ctx, params)
	}
	return []entities.PutawaySuggestion{}, 0, nil
}

func (m *UseCaseMock) CreateReplenishmentTask(ctx context.Context, dto request.CreateReplenishmentTaskDTO) (*entities.ReplenishmentTask, error) {
	if m.CreateReplenishmentTaskFn != nil {
		return m.CreateReplenishmentTaskFn(ctx, dto)
	}
	return &entities.ReplenishmentTask{}, nil
}

func (m *UseCaseMock) ListReplenishmentTasks(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error) {
	if m.ListReplenishmentTasksFn != nil {
		return m.ListReplenishmentTasksFn(ctx, params)
	}
	return []entities.ReplenishmentTask{}, 0, nil
}

func (m *UseCaseMock) AssignReplenishment(ctx context.Context, dto request.AssignReplenishmentDTO) (*entities.ReplenishmentTask, error) {
	if m.AssignReplenishmentFn != nil {
		return m.AssignReplenishmentFn(ctx, dto)
	}
	return &entities.ReplenishmentTask{}, nil
}

func (m *UseCaseMock) CompleteReplenishment(ctx context.Context, dto request.CompleteReplenishmentDTO) (*entities.ReplenishmentTask, error) {
	if m.CompleteReplenishmentFn != nil {
		return m.CompleteReplenishmentFn(ctx, dto)
	}
	return &entities.ReplenishmentTask{}, nil
}

func (m *UseCaseMock) CancelReplenishment(ctx context.Context, businessID, taskID uint, reason string) (*entities.ReplenishmentTask, error) {
	if m.CancelReplenishmentFn != nil {
		return m.CancelReplenishmentFn(ctx, businessID, taskID, reason)
	}
	return &entities.ReplenishmentTask{}, nil
}

func (m *UseCaseMock) DetectReplenishmentNeeds(ctx context.Context, businessID uint) (*response.ReplenishmentDetectResult, error) {
	if m.DetectReplenishmentNeedsFn != nil {
		return m.DetectReplenishmentNeedsFn(ctx, businessID)
	}
	return &response.ReplenishmentDetectResult{}, nil
}

func (m *UseCaseMock) CreateCrossDockLink(ctx context.Context, dto request.CreateCrossDockLinkDTO) (*entities.CrossDockLink, error) {
	if m.CreateCrossDockLinkFn != nil {
		return m.CreateCrossDockLinkFn(ctx, dto)
	}
	return &entities.CrossDockLink{}, nil
}

func (m *UseCaseMock) ListCrossDockLinks(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error) {
	if m.ListCrossDockLinksFn != nil {
		return m.ListCrossDockLinksFn(ctx, params)
	}
	return []entities.CrossDockLink{}, 0, nil
}

func (m *UseCaseMock) ExecuteCrossDock(ctx context.Context, dto request.ExecuteCrossDockDTO) (*entities.CrossDockLink, error) {
	if m.ExecuteCrossDockFn != nil {
		return m.ExecuteCrossDockFn(ctx, dto)
	}
	return &entities.CrossDockLink{}, nil
}

func (m *UseCaseMock) RunSlotting(ctx context.Context, dto request.RunSlottingDTO) (*response.SlottingRunResult, error) {
	if m.RunSlottingFn != nil {
		return m.RunSlottingFn(ctx, dto)
	}
	return &response.SlottingRunResult{}, nil
}

func (m *UseCaseMock) ListVelocities(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error) {
	if m.ListVelocitiesFn != nil {
		return m.ListVelocitiesFn(ctx, params)
	}
	return []entities.ProductVelocity{}, nil
}
