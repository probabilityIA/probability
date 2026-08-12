package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

func (uc *useCase) FindingItems(ctx context.Context, query domain.FindingItemsQuery) (*domain.FindingItemsPage, error) {
	if query.BusinessID == 0 || query.Code == "" {
		return &domain.FindingItemsPage{Items: []domain.FindingItem{}, Page: 1, PageSize: 50}, nil
	}
	query.Normalize()
	return uc.repo.FindingItems(ctx, query)
}
