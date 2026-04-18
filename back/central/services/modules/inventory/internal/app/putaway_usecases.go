package app

import (
	"context"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) CreatePutawayRule(ctx context.Context, dto request.CreatePutawayRuleDTO) (*entities.PutawayRule, error) {
	strategy := dto.Strategy
	if strategy == "" {
		strategy = "nearest_empty"
	}
	rule := &entities.PutawayRule{
		BusinessID:   dto.BusinessID,
		ProductID:    dto.ProductID,
		CategoryID:   dto.CategoryID,
		TargetZoneID: dto.TargetZoneID,
		Priority:     dto.Priority,
		Strategy:     strategy,
		IsActive:     dto.IsActive,
	}
	return uc.repo.CreatePutawayRule(ctx, rule)
}

func (uc *useCase) ListPutawayRules(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListPutawayRules(ctx, params)
}

func (uc *useCase) UpdatePutawayRule(ctx context.Context, dto request.UpdatePutawayRuleDTO) (*entities.PutawayRule, error) {
	existing, err := uc.repo.GetPutawayRuleByID(ctx, dto.BusinessID, dto.ID)
	if err != nil {
		return nil, err
	}
	if dto.ProductID != nil {
		existing.ProductID = dto.ProductID
	}
	if dto.CategoryID != nil {
		existing.CategoryID = dto.CategoryID
	}
	if dto.TargetZoneID != nil {
		existing.TargetZoneID = *dto.TargetZoneID
	}
	if dto.Priority != nil {
		existing.Priority = *dto.Priority
	}
	if dto.Strategy != "" {
		existing.Strategy = dto.Strategy
	}
	if dto.IsActive != nil {
		existing.IsActive = *dto.IsActive
	}
	return uc.repo.UpdatePutawayRule(ctx, existing)
}

func (uc *useCase) DeletePutawayRule(ctx context.Context, businessID, ruleID uint) error {
	return uc.repo.DeletePutawayRule(ctx, businessID, ruleID)
}

func (uc *useCase) SuggestPutaway(ctx context.Context, dto request.PutawaySuggestDTO) (*response.PutawaySuggestResult, error) {
	result := &response.PutawaySuggestResult{}
	for _, item := range dto.Items {
		rule, err := uc.repo.FindApplicableRule(ctx, dto.BusinessID, item.ProductID)
		if err != nil {
			result.UnresolvedItems = append(result.UnresolvedItems, item.ProductID)
			continue
		}
		locID, err := uc.repo.PickLocationInZone(ctx, rule.TargetZoneID)
		if err != nil {
			result.UnresolvedItems = append(result.UnresolvedItems, item.ProductID)
			continue
		}
		sug := &entities.PutawaySuggestion{
			BusinessID:            dto.BusinessID,
			ProductID:             item.ProductID,
			RecommendedLocationID: locID,
			Quantity:              item.Quantity,
			Status:                "pending",
			RuleID:                &rule.ID,
			Reason:                rule.Strategy,
		}
		created, err := uc.repo.CreatePutawaySuggestion(ctx, sug)
		if err != nil {
			result.UnresolvedItems = append(result.UnresolvedItems, item.ProductID)
			continue
		}
		result.Suggestions = append(result.Suggestions, *created)
	}
	return result, nil
}

func (uc *useCase) ConfirmPutaway(ctx context.Context, dto request.ConfirmPutawayDTO) (*entities.PutawaySuggestion, error) {
	existing, err := uc.repo.GetPutawaySuggestionByID(ctx, dto.BusinessID, dto.SuggestionID)
	if err != nil {
		return nil, err
	}
	if existing.Status == "confirmed" {
		return nil, domainerrors.ErrPutawayAlreadyConfirmed
	}
	now := time.Now()
	existing.Status = "confirmed"
	existing.ActualLocationID = &dto.ActualLocationID
	existing.ConfirmedAt = &now
	existing.ConfirmedByID = dto.UserID
	return uc.repo.UpdatePutawaySuggestion(ctx, existing)
}

func (uc *useCase) ListPutawaySuggestions(ctx context.Context, params dtos.ListPutawaySuggestionsParams) ([]entities.PutawaySuggestion, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.ListPutawaySuggestions(ctx, params)
}

func (uc *useCase) NoPutawayRuleFound(err error) bool {
	return errors.Is(err, domainerrors.ErrNoPutawayRule)
}
