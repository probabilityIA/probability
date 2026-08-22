package core

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (j *JumpsellerCore) SupportsOrdersCompare() bool {
	return true
}

func (j *JumpsellerCore) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	return j.useCase.ListChannelOrders(ctx, integrationID, filters)
}

func (j *JumpsellerCore) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	return j.useCase.ImportChannelOrders(ctx, integrationID, externalIDs)
}
