package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (uc *meliUseCase) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}

	sellerID, err := extractSellerID(integration)
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("ensuring valid token: %w", err)
	}

	cli := uc.clientFor(ctx, integration)
	limit := filters.Limit
	if limit <= 0 {
		limit = 50
	}

	params := &domain.GetOrdersParams{
		DateFrom: filters.From,
		DateTo:   filters.To,
		Offset:   0,
		Limit:    50,
		Sort:     "date_desc",
	}

	rows := make([]orderscompare.ChannelOrder, 0, limit)
	seenPacks := make(map[int64]bool)

	for {
		result, _, err := cli.GetOrders(ctx, accessToken, sellerID, params)
		if err != nil {
			return nil, fmt.Errorf("obteniendo ordenes de MercadoLibre: %w", err)
		}
		if len(result.Orders) == 0 {
			break
		}

		for i := range result.Orders {
			order := result.Orders[i]
			if order.PackID != nil && *order.PackID > 0 {
				if seenPacks[*order.PackID] {
					continue
				}
				seenPacks[*order.PackID] = true
			}

			rows = append(rows, orderscompare.ChannelOrder{
				ExternalID:   strconv.FormatInt(order.ID, 10),
				Number:       strconv.FormatInt(order.ID, 10),
				CustomerName: meliBuyerName(&order),
				Status:       mapper.MapMeliStatus(order.Status),
				RawStatus:    order.Status,
				Total:        order.TotalAmount,
				Currency:     order.CurrencyID,
				Items:        len(order.OrderItems),
				CreatedAt:    order.DateCreated,
			})
			if len(rows) >= limit {
				return rows, nil
			}
		}

		params.Offset += params.Limit
		if params.Offset >= result.Total {
			break
		}
	}

	return rows, nil
}

func (uc *meliUseCase) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	result := orderscompare.ImportResult{
		Queued:   []string{},
		Failed:   map[string]string{},
		NotFound: []string{},
	}

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return result, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return result, domain.ErrIntegrationNotFound
	}

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return result, fmt.Errorf("ensuring valid token: %w", err)
	}

	cli := uc.clientFor(ctx, integration)

	for _, externalID := range externalIDs {
		orderID, err := strconv.ParseInt(externalID, 10, 64)
		if err != nil {
			result.Failed[externalID] = "external_id invalido para MercadoLibre"
			continue
		}

		order, rawJSON, err := cli.GetOrder(ctx, accessToken, orderID)
		if err != nil {
			uc.logger.Warn(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", externalID).
				Msg("No se pudo leer la orden de MercadoLibre para importarla")
			result.NotFound = append(result.NotFound, externalID)
			continue
		}

		var shippingDetail *domain.MeliShippingDetail
		var shipmentRaw []byte
		if order.Shipping != nil && order.Shipping.ID > 0 {
			if detail, raw, shErr := cli.GetShipmentDetail(ctx, accessToken, order.Shipping.ID); shErr == nil {
				shippingDetail = detail
				shipmentRaw = raw
			}
		}

		dto := mapper.MapMeliOrderToProbability(order, shippingDetail, rawJSON, shipmentRaw)
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

func meliBuyerName(order *domain.MeliOrder) string {
	name := strings.TrimSpace(order.Buyer.FirstName + " " + order.Buyer.LastName)
	if name != "" {
		return name
	}
	return order.Buyer.Nickname
}
