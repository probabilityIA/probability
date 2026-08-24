package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

func (uc *tiendanubeUseCase) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error) {
	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return nil, err
	}

	orders, err := uc.client.GetOrders(ctx, cred, toTiendanubeFilters(filters))
	if err != nil {
		return nil, fmt.Errorf("obteniendo ordenes de Tiendanube: %w", err)
	}

	rows := make([]orderscompare.ChannelOrder, 0, len(orders))
	for i := range orders {
		order := orders[i]
		rows = append(rows, orderscompare.ChannelOrder{
			ExternalID:   order.ID,
			Number:       order.Number,
			CustomerName: order.ContactName,
			Status:       mapper.MapOrderStatus(order.Status, order.PaymentStatus, order.ShippingStatus),
			RawStatus:    order.Status,
			Total:        order.Total,
			Currency:     order.Currency,
			Items:        len(order.Items),
			CreatedAt:    parseTiendanubeDate(order.CreatedAt),
		})
	}
	return rows, nil
}

func (uc *tiendanubeUseCase) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error) {
	result := orderscompare.ImportResult{
		Queued:   []string{},
		Failed:   map[string]string{},
		NotFound: []string{},
	}

	integration, err := uc.fetchIntegration(ctx, integrationID)
	if err != nil {
		return result, err
	}
	cred, err := uc.buildCredential(ctx, integrationID, integration)
	if err != nil {
		return result, err
	}

	for _, externalID := range externalIDs {
		order, rawJSON, err := uc.client.GetOrder(ctx, cred, externalID)
		if err != nil {
			uc.logger.Warn(ctx).Err(err).
				Str("integration_id", integrationID).
				Str("order_id", externalID).
				Msg("No se pudo leer la orden de Tiendanube para importarla")
			result.NotFound = append(result.NotFound, externalID)
			continue
		}

		skip, _ := orderscompare.SkipsInventory(mapper.MapOrderStatus(order.Status, order.PaymentStatus, order.ShippingStatus))
		if err := uc.publishOrderWithInventoryPolicy(ctx, integration, order, rawJSON, skip); err != nil {
			result.Failed[externalID] = err.Error()
			continue
		}
		result.Queued = append(result.Queued, externalID)
	}

	return result, nil
}

func (uc *tiendanubeUseCase) publishOrderWithInventoryPolicy(ctx context.Context, integration *domain.Integration, order *domain.TiendanubeOrder, rawJSON []byte, skipInventory bool) error {
	dto := mapper.MapTiendanubeOrderToProbability(order, rawJSON)
	dto.IntegrationID = integration.ID
	dto.BusinessID = integration.BusinessID
	dto.SkipInventory = skipInventory

	if err := uc.publisher.Publish(ctx, dto); err != nil {
		return fmt.Errorf("publicando la orden %s de Tiendanube: %w", order.ID, err)
	}
	return nil
}

func toTiendanubeFilters(filters orderscompare.ChannelFilters) domain.OrderFilters {
	out := domain.OrderFilters{Limit: filters.Limit}
	if filters.From != nil {
		out.CreatedAtMin = filters.From.Format(time.RFC3339)
	}
	if filters.To != nil {
		out.CreatedAtMax = filters.To.Format(time.RFC3339)
	}
	return out
}

func parseTiendanubeDate(value string) time.Time {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return time.Time{}
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05-0700", "2006-01-02 15:04:05"} {
		if parsed, err := time.Parse(layout, clean); err == nil {
			return parsed
		}
	}
	return time.Time{}
}
