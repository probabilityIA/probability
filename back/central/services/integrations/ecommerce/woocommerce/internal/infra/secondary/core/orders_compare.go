package core

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (w *WooCommerceCore) SupportsOrdersCompare() bool {
	return true
}

func (w *WooCommerceCore) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	return w.useCase.ListChannelOrders(ctx, integrationID, filters)
}

func (w *WooCommerceCore) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	return w.useCase.ImportChannelOrders(ctx, integrationID, externalIDs)
}
