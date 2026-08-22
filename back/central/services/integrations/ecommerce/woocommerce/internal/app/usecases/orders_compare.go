package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

type wooAccess struct {
	integration    *domain.Integration
	storeURL       string
	consumerKey    string
	consumerSecret string
}

func (uc *wooCommerceUseCase) resolveAccess(ctx context.Context, integrationID string) (*wooAccess, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}

	storeURL, err := extractString(integration.Config, "store_url")
	if err != nil {
		return nil, domain.ErrMissingStoreURL
	}
	storeURL = resolveEffectiveStoreURL(integration, storeURL)

	consumerKey, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_key")
	if err != nil {
		return nil, fmt.Errorf("decrypting consumer_key: %w", err)
	}
	consumerSecret, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_secret")
	if err != nil {
		return nil, fmt.Errorf("decrypting consumer_secret: %w", err)
	}

	return &wooAccess{
		integration:    integration,
		storeURL:       storeURL,
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
	}, nil
}

func (uc *wooCommerceUseCase) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	access, err := uc.resolveAccess(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	perPage := filters.Limit
	if perPage <= 0 || perPage > 100 {
		perPage = 100
	}

	params := &domain.GetOrdersParams{
		After:   filters.From,
		Before:  filters.To,
		Page:    1,
		PerPage: perPage,
		OrderBy: "date",
		Order:   "desc",
	}

	rows := make([]orderscompare.ChannelOrder, 0, perPage)
	for {
		result, _, err := uc.client.GetOrders(ctx, access.storeURL, access.consumerKey, access.consumerSecret, params)
		if err != nil {
			return nil, fmt.Errorf("obteniendo ordenes de WooCommerce: %w", err)
		}
		if len(result.Orders) == 0 {
			break
		}

		for i := range result.Orders {
			order := result.Orders[i]
			rows = append(rows, orderscompare.ChannelOrder{
				ExternalID:   strconv.FormatInt(order.ID, 10),
				Number:       order.Number,
				CustomerName: customerNameOf(&order),
				Status:       mapper.MapWooStatus(order.Status),
				RawStatus:    order.Status,
				Total:        parseWooTotal(order.Total),
				Currency:     order.Currency,
				Items:        len(order.LineItems),
				CreatedAt:    order.DateCreated,
			})
			if filters.Limit > 0 && len(rows) >= filters.Limit {
				return rows, nil
			}
		}

		params.Page++
		if params.Page > 20 {
			break
		}
	}

	return rows, nil
}

func (uc *wooCommerceUseCase) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	result := orderscompare.ImportResult{
		Queued:   []string{},
		Failed:   map[string]string{},
		NotFound: []string{},
	}

	access, err := uc.resolveAccess(ctx, integrationID)
	if err != nil {
		return result, err
	}

	for _, externalID := range externalIDs {
		orderID, err := strconv.ParseInt(externalID, 10, 64)
		if err != nil {
			result.Failed[externalID] = "external_id invalido para WooCommerce"
			continue
		}

		order, rawJSON, err := uc.client.GetOrder(ctx, access.storeURL, access.consumerKey, access.consumerSecret, orderID)
		if err != nil {
			uc.logger.Warn(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", externalID).
				Msg("No se pudo leer la orden de WooCommerce para importarla")
			result.NotFound = append(result.NotFound, externalID)
			continue
		}

		dto := mapper.MapWooOrderToProbability(order, rawJSON)
		dto.IntegrationID = access.integration.ID
		dto.BusinessID = access.integration.BusinessID
		skip, _ := orderscompare.SkipsInventory(dto.Status)
		dto.SkipInventory = skip

		if err := uc.publisher.Publish(ctx, dto); err != nil {
			result.Failed[externalID] = err.Error()
			continue
		}
		result.Queued = append(result.Queued, externalID)
	}

	return result, nil
}

func customerNameOf(order *domain.WooCommerceOrder) string {
	name := order.Billing.FirstName + " " + order.Billing.LastName
	if len(name) > 1 {
		return name
	}
	return order.Shipping.FirstName + " " + order.Shipping.LastName
}

func parseWooTotal(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return parsed
}
