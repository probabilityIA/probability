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
	posicionPorClave := make(map[string]int, limit)
	packsConsultados := make(map[int64]bool, limit)

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
			clave := uc.claveComparativo(ctx, cli, accessToken, &order, packsConsultados)

			if pos, ok := posicionPorClave[clave]; ok {
				rows[pos].Total += order.TotalAmount
				rows[pos].Items += len(order.OrderItems)
				continue
			}

			rows = append(rows, orderscompare.ChannelOrder{
				ExternalID:        clave,
				Number:            clave,
				CustomerName:      meliBuyerName(&order),
				Status:            mapper.MapMeliStatus(order.Status),
				RawStatus:         order.Status,
				FulfillmentStatus: meliFulfillment(&order),
				Total:             order.TotalAmount,
				Currency:          order.CurrencyID,
				Items:             len(order.OrderItems),
				CreatedAt:         order.DateCreated,
			})
			posicionPorClave[clave] = len(rows) - 1

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
		id, err := strconv.ParseInt(externalID, 10, 64)
		if err != nil {
			result.Failed[externalID] = "external_id invalido para MercadoLibre"
			continue
		}

		order, rawJSON, err := uc.resolveOrderOrPack(ctx, cli, accessToken, id)
		if err != nil || order == nil {
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

		uc.enrichBillingInfo(ctx, integrationID, accessToken, order)

		dto := mapper.MapMeliOrderToProbability(order, shippingDetail, rawJSON, shipmentRaw)
		dto.IntegrationID = integration.ID
		dto.BusinessID = integration.BusinessID
		skip, _ := orderscompare.SkipsInventoryFor(dto.Status, meliFulfillmentDetallado(order, shippingDetail))
		dto.SkipInventory = skip

		if err := uc.publisher.Publish(ctx, dto); err != nil {
			result.Failed[externalID] = err.Error()
			continue
		}
		result.Queued = append(result.Queued, externalID)
	}

	return result, nil
}

func (uc *meliUseCase) claveComparativo(ctx context.Context, cli domain.IMeliClient, accessToken string, order *domain.MeliOrder, consultados map[int64]bool) string {
	if order.PackID == nil || *order.PackID <= 0 {
		return strconv.FormatInt(order.ID, 10)
	}

	packID := *order.PackID
	if consultados[packID] {
		return strconv.FormatInt(packID, 10)
	}

	pack, err := cli.GetPack(ctx, accessToken, packID)
	if err != nil {
		uc.logger.Warn(ctx).Err(err).Int64("pack_id", packID).
			Msg("No se pudo leer el pack de MercadoLibre, se compara por orden")
		return strconv.FormatInt(order.ID, 10)
	}
	if pack == nil || len(pack.OrderIDs) <= 1 {
		return strconv.FormatInt(order.ID, 10)
	}

	consultados[packID] = true
	return strconv.FormatInt(packID, 10)
}

func (uc *meliUseCase) resolveOrderOrPack(ctx context.Context, cli domain.IMeliClient, accessToken string, id int64) (*domain.MeliOrder, []byte, error) {
	order, rawJSON, err := cli.GetOrder(ctx, accessToken, id)
	if err == nil && order != nil {
		if order.PackID != nil && *order.PackID > 0 {
			if merged, perr := uc.consolidatePack(ctx, accessToken, *order.PackID, order); perr == nil && merged != nil {
				return merged, rawJSON, nil
			}
		}
		return order, rawJSON, nil
	}

	pack, perr := cli.GetPack(ctx, accessToken, id)
	if perr != nil || pack == nil || len(pack.OrderIDs) == 0 {
		return nil, nil, err
	}

	primary, primaryRaw, gerr := cli.GetOrder(ctx, accessToken, pack.OrderIDs[0])
	if gerr != nil || primary == nil {
		return nil, nil, gerr
	}

	merged, merr := uc.consolidatePack(ctx, accessToken, id, primary)
	if merr != nil || merged == nil {
		return primary, primaryRaw, nil
	}
	return merged, primaryRaw, nil
}

func meliFulfillment(order *domain.MeliOrder) string {
	for _, tag := range order.Tags {
		switch strings.ToLower(strings.TrimSpace(tag)) {
		case "delivered":
			return "delivered"
		case "not_delivered":
			return ""
		}
	}
	return ""
}

func meliFulfillmentDetallado(order *domain.MeliOrder, shipping *domain.MeliShippingDetail) string {
	if fulfillment := meliFulfillment(order); fulfillment != "" {
		return fulfillment
	}
	if shipping == nil {
		return ""
	}
	return shipping.Status
}

func meliBuyerName(order *domain.MeliOrder) string {
	name := strings.TrimSpace(order.Buyer.FirstName + " " + order.Buyer.LastName)
	if name != "" {
		return name
	}
	return order.Buyer.Nickname
}
