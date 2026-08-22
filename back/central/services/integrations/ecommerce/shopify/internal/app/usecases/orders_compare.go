package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/utils"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

type shopifyAccess struct {
	integration *domain.Integration
	storeDomain string
	accessToken string
}

func (uc *SyncOrdersUseCase) resolveAccess(ctx context.Context, integrationID string) (*shopifyAccess, error) {
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	config, err := utils.NormalizeConfig(integration.Config, integration.Name)
	if err != nil {
		return nil, err
	}

	storeDomain, err := utils.ExtractStoreName(config, integration.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to extract store name: %w", err)
	}
	storeDomain = utils.ResolveEffectiveStoreDomain(integration, storeDomain)

	accessToken, err := utils.GetAccessToken(ctx, uc.integrationService, integrationID)
	if err != nil {
		return nil, err
	}

	return &shopifyAccess{integration: integration, storeDomain: storeDomain, accessToken: accessToken}, nil
}

func (uc *SyncOrdersUseCase) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	access, err := uc.resolveAccess(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	limit := filters.Limit
	if limit <= 0 || limit > 250 {
		limit = 250
	}

	params := &domain.GetOrdersParams{
		Limit:             limit,
		CreatedAtMin:      filters.From,
		CreatedAtMax:      filters.To,
		Status:            domain.OrderStatusAny,
		FinancialStatus:   domain.FinancialStatusAny,
		FulfillmentStatus: domain.FulfillmentStatusAny,
	}

	orders, _, err := uc.shopifyClient.GetOrders(ctx, access.storeDomain, access.accessToken, params)
	if err != nil {
		return nil, fmt.Errorf("obteniendo ordenes de Shopify: %w", err)
	}

	rows := make([]orderscompare.ChannelOrder, 0, len(orders))
	for i := range orders {
		order := orders[i]
		rows = append(rows, orderscompare.ChannelOrder{
			ExternalID:   order.ExternalID,
			Number:       order.OrderNumber,
			CustomerName: shopifyCustomerName(&order),
			Status:       order.Status,
			RawStatus:    order.OriginalStatus,
			Total:        order.TotalAmount,
			Currency:     order.Currency,
			Items:        len(order.Items),
			CreatedAt:    order.OccurredAt,
			URL:          order.OrderStatusURL,
		})
	}
	return rows, nil
}

func (uc *SyncOrdersUseCase) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	result := orderscompare.ImportResult{
		Queued:   []string{},
		Failed:   map[string]string{},
		NotFound: []string{},
	}

	access, err := uc.resolveAccess(ctx, integrationID)
	if err != nil {
		return result, err
	}
	if access.integration.BusinessID == nil {
		return result, domain.ErrBusinessIDMissing
	}

	for _, externalID := range externalIDs {
		order, err := uc.shopifyClient.GetOrder(ctx, access.storeDomain, access.accessToken, externalID)
		if err != nil || order == nil {
			uc.log.Warn(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", externalID).
				Msg("No se pudo leer la orden de Shopify para importarla")
			result.NotFound = append(result.NotFound, externalID)
			continue
		}

		order.IntegrationID = access.integration.ID
		order.IntegrationType = "shopify"
		order.BusinessID = access.integration.BusinessID

		dto := mapper.MapShopifyOrderToProbability(order)
		mapper.EnrichOrderWithDetails(dto, order.RawData)
		skip, _ := orderscompare.SkipsInventory(dto.Status)
		dto.SkipInventory = skip

		if err := uc.orderPublisher.Publish(ctx, dto); err != nil {
			result.Failed[externalID] = err.Error()
			continue
		}
		result.Queued = append(result.Queued, externalID)
	}

	return result, nil
}

func shopifyCustomerName(order *domain.ShopifyOrder) string {
	name := strings.TrimSpace(order.Customer.Name)
	if name != "" {
		return name
	}
	return strings.TrimSpace(order.Customer.Email)
}
