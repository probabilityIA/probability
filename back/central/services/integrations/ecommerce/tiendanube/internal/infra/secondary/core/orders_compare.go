package core

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (t *TiendanubeCore) SupportsOrdersCompare() bool {
	return true
}

func (t *TiendanubeCore) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	return t.useCase.ListChannelOrders(ctx, integrationID, filters)
}

func (t *TiendanubeCore) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	return t.useCase.ImportChannelOrders(ctx, integrationID, externalIDs)
}
