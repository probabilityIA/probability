package core

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (v *VTEXCore) SupportsOrdersCompare() bool {
	return true
}

func (v *VTEXCore) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	return v.useCase.ListChannelOrders(ctx, integrationID, filters)
}

func (v *VTEXCore) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	return v.useCase.ImportChannelOrders(ctx, integrationID, externalIDs)
}
