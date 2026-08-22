package core

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (m *MeliCore) SupportsOrdersCompare() bool {
	return true
}

func (m *MeliCore) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	return m.useCase.ListChannelOrders(ctx, integrationID, filters)
}

func (m *MeliCore) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	return m.useCase.ImportChannelOrders(ctx, integrationID, externalIDs)
}
