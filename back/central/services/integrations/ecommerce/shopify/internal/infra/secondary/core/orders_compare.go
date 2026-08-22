package core

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (s *ShopifyCore) SupportsOrdersCompare() bool {
	return true
}

func (s *ShopifyCore) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	return s.useCase.ListChannelOrders(ctx, integrationID, filters)
}

func (s *ShopifyCore) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	return s.useCase.ImportChannelOrders(ctx, integrationID, externalIDs)
}
