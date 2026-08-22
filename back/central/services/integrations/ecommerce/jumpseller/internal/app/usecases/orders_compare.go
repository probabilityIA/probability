package usecases

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (uc *jumpsellerUseCase) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	_, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	limit := filters.Limit
	if limit <= 0 {
		limit = ordersPageSize
	}

	raw := map[string]interface{}{}
	if filters.From != nil {
		raw["created_at_min"] = filters.From.Format(time.RFC3339)
	}
	if filters.To != nil {
		raw["created_at_max"] = filters.To.Format(time.RFC3339)
	}
	params, err := buildQueryParams(raw)
	if err != nil {
		return nil, err
	}
	params.PerPage = ordersPageSize

	rows := make([]orderscompare.ChannelOrder, 0, limit)
	for page := 1; page <= maxOrderPages; page++ {
		params.Page = page

		result, _, err := uc.client.GetOrders(ctx, cred, params)
		if err != nil {
			return nil, fmt.Errorf("obteniendo ordenes de Jumpseller: %w", err)
		}
		if len(result.Orders) == 0 {
			break
		}

		for i := range result.Orders {
			order := result.Orders[i]
			rows = append(rows, orderscompare.ChannelOrder{
				ExternalID:   strconv.FormatInt(order.ID, 10),
				Number:       strconv.FormatInt(order.ID, 10),
				CustomerName: order.Customer.Name,
				Status:       mapper.MapJumpsellerOrderStatus(order.Status, order.StatusEnum),
				RawStatus:    order.Status,
				Total:        order.Total,
				Currency:     order.Currency,
				Items:        len(order.Products),
				CreatedAt:    order.CreatedAt,
			})
			if len(rows) >= limit {
				return rows, nil
			}
		}
	}

	return rows, nil
}

func (uc *jumpsellerUseCase) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	result := orderscompare.ImportResult{
		Queued:   []string{},
		Failed:   map[string]string{},
		NotFound: []string{},
	}

	integration, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return result, err
	}

	for _, externalID := range externalIDs {
		orderID, err := strconv.ParseInt(externalID, 10, 64)
		if err != nil {
			result.Failed[externalID] = "external_id invalido para Jumpseller"
			continue
		}

		order, rawJSON, err := uc.client.GetOrder(ctx, cred, orderID)
		if err != nil {
			uc.logger.Warn(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", externalID).
				Msg("No se pudo leer la orden de Jumpseller para importarla")
			result.NotFound = append(result.NotFound, externalID)
			continue
		}

		dto := mapper.MapJumpsellerOrderToProbability(order, rawJSON)
		dto.IntegrationID = integration.ID
		dto.BusinessID = integration.BusinessID
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
