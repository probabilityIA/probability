package channels

import (
	"context"
	"fmt"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

type registry struct {
	core integrationcore.IIntegrationCore
}

func New(core integrationcore.IIntegrationCore) ports.IChannelRegistry {
	return &registry{core: core}
}

func (r *registry) Supports(integrationTypeID uint) bool {
	contract, ok := r.core.GetRegisteredIntegration(int(integrationTypeID))
	if !ok || contract == nil {
		return false
	}
	return contract.SupportsOrdersCompare()
}

func (r *registry) ListOrders(ctx context.Context, integrationTypeID uint, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	contract, ok := r.core.GetRegisteredIntegration(int(integrationTypeID))
	if !ok || contract == nil {
		return nil, fmt.Errorf("no hay un canal registrado para el tipo %d", integrationTypeID)
	}

	return contract.ListChannelOrders(ctx, integrationID, filters)
}

func (r *registry) ImportOrders(ctx context.Context, integrationTypeID uint, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	contract, ok := r.core.GetRegisteredIntegration(int(integrationTypeID))
	if !ok || contract == nil {
		return orderscompare.ImportResult{}, fmt.Errorf("no hay un canal registrado para el tipo %d", integrationTypeID)
	}

	return contract.ImportChannelOrders(ctx, integrationID, externalIDs)
}
