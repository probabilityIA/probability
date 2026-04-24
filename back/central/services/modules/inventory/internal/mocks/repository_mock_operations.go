package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

func (m *RepositoryMock) CreatePutawayRule(ctx context.Context, rule *entities.PutawayRule) (*entities.PutawayRule, error) {
	if m.CreatePutawayRuleFn != nil {
		return m.CreatePutawayRuleFn(ctx, rule)
	}
	return rule, nil
}

func (m *RepositoryMock) ListPutawayRules(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error) {
	if m.ListPutawayRulesFn != nil {
		return m.ListPutawayRulesFn(ctx, params)
	}
	return []entities.PutawayRule{}, 0, nil
}

func (m *RepositoryMock) GetPutawayRuleByID(ctx context.Context, businessID, ruleID uint) (*entities.PutawayRule, error) {
	if m.GetPutawayRuleByIDFn != nil {
		return m.GetPutawayRuleByIDFn(ctx, businessID, ruleID)
	}
	return &entities.PutawayRule{ID: ruleID, BusinessID: businessID}, nil
}

func (m *RepositoryMock) UpdatePutawayRule(ctx context.Context, rule *entities.PutawayRule) (*entities.PutawayRule, error) {
	if m.UpdatePutawayRuleFn != nil {
		return m.UpdatePutawayRuleFn(ctx, rule)
	}
	return rule, nil
}

func (m *RepositoryMock) DeletePutawayRule(ctx context.Context, businessID, ruleID uint) error {
	if m.DeletePutawayRuleFn != nil {
		return m.DeletePutawayRuleFn(ctx, businessID, ruleID)
	}
	return nil
}

func (m *RepositoryMock) FindApplicableRule(ctx context.Context, businessID uint, productID string) (*entities.PutawayRule, error) {
	if m.FindApplicableRuleFn != nil {
		return m.FindApplicableRuleFn(ctx, businessID, productID)
	}
	return &entities.PutawayRule{}, nil
}

func (m *RepositoryMock) PickLocationInZone(ctx context.Context, zoneID uint) (uint, error) {
	if m.PickLocationInZoneFn != nil {
		return m.PickLocationInZoneFn(ctx, zoneID)
	}
	return 0, nil
}

func (m *RepositoryMock) CreatePutawaySuggestion(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
	if m.CreatePutawaySuggestionFn != nil {
		return m.CreatePutawaySuggestionFn(ctx, s)
	}
	return s, nil
}

func (m *RepositoryMock) GetPutawaySuggestionByID(ctx context.Context, businessID, id uint) (*entities.PutawaySuggestion, error) {
	if m.GetPutawaySuggestionByIDFn != nil {
		return m.GetPutawaySuggestionByIDFn(ctx, businessID, id)
	}
	return &entities.PutawaySuggestion{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListPutawaySuggestions(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error) {
	if m.ListPutawaySuggestionsFn != nil {
		return m.ListPutawaySuggestionsFn(ctx, params)
	}
	return []entities.PutawaySuggestion{}, 0, nil
}

func (m *RepositoryMock) UpdatePutawaySuggestion(ctx context.Context, s *entities.PutawaySuggestion) (*entities.PutawaySuggestion, error) {
	if m.UpdatePutawaySuggestionFn != nil {
		return m.UpdatePutawaySuggestionFn(ctx, s)
	}
	return s, nil
}

func (m *RepositoryMock) CreateReplenishmentTask(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
	if m.CreateReplenishmentTaskFn != nil {
		return m.CreateReplenishmentTaskFn(ctx, t)
	}
	return t, nil
}

func (m *RepositoryMock) GetReplenishmentTaskByID(ctx context.Context, businessID, id uint) (*entities.ReplenishmentTask, error) {
	if m.GetReplenishmentTaskByIDFn != nil {
		return m.GetReplenishmentTaskByIDFn(ctx, businessID, id)
	}
	return &entities.ReplenishmentTask{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListReplenishmentTasks(ctx context.Context, params dtos.ListReplenishmentTasksParams) ([]entities.ReplenishmentTask, int64, error) {
	if m.ListReplenishmentTasksFn != nil {
		return m.ListReplenishmentTasksFn(ctx, params)
	}
	return []entities.ReplenishmentTask{}, 0, nil
}

func (m *RepositoryMock) UpdateReplenishmentTask(ctx context.Context, t *entities.ReplenishmentTask) (*entities.ReplenishmentTask, error) {
	if m.UpdateReplenishmentTaskFn != nil {
		return m.UpdateReplenishmentTaskFn(ctx, t)
	}
	return t, nil
}

func (m *RepositoryMock) DetectReplenishmentCandidates(ctx context.Context, businessID uint) ([]entities.ReplenishmentTask, error) {
	if m.DetectReplenishmentCandidatesFn != nil {
		return m.DetectReplenishmentCandidatesFn(ctx, businessID)
	}
	return []entities.ReplenishmentTask{}, nil
}

func (m *RepositoryMock) CreateCrossDockLink(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error) {
	if m.CreateCrossDockLinkFn != nil {
		return m.CreateCrossDockLinkFn(ctx, l)
	}
	return l, nil
}

func (m *RepositoryMock) GetCrossDockLinkByID(ctx context.Context, businessID, id uint) (*entities.CrossDockLink, error) {
	if m.GetCrossDockLinkByIDFn != nil {
		return m.GetCrossDockLinkByIDFn(ctx, businessID, id)
	}
	return &entities.CrossDockLink{ID: id, BusinessID: businessID}, nil
}

func (m *RepositoryMock) ListCrossDockLinks(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error) {
	if m.ListCrossDockLinksFn != nil {
		return m.ListCrossDockLinksFn(ctx, params)
	}
	return []entities.CrossDockLink{}, 0, nil
}

func (m *RepositoryMock) UpdateCrossDockLink(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error) {
	if m.UpdateCrossDockLinkFn != nil {
		return m.UpdateCrossDockLinkFn(ctx, l)
	}
	return l, nil
}

func (m *RepositoryMock) ComputeVelocities(ctx context.Context, businessID, warehouseID uint, period string) error {
	if m.ComputeVelocitiesFn != nil {
		return m.ComputeVelocitiesFn(ctx, businessID, warehouseID, period)
	}
	return nil
}

func (m *RepositoryMock) ListVelocities(ctx context.Context, params dtos.ListVelocityParams) ([]entities.ProductVelocity, error) {
	if m.ListVelocitiesFn != nil {
		return m.ListVelocitiesFn(ctx, params)
	}
	return []entities.ProductVelocity{}, nil
}
