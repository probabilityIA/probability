package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func (uc *meliUseCase) emitOrderSyncEvent(ctx context.Context, integration *domain.Integration, eventType string, data map[string]interface{}) {
	if uc.rabbit == nil {
		return
	}
	var bID uint
	if integration.BusinessID != nil {
		bID = *integration.BusinessID
	}
	_ = rabbitmq.PublishEvent(ctx, uc.rabbit, rabbitmq.EventEnvelope{
		Type:          eventType,
		Category:      "integration",
		BusinessID:    bID,
		IntegrationID: integration.ID,
		Data:          data,
	})
}

func (uc *meliUseCase) SyncOrders(ctx context.Context, integrationID string) error {
	after := time.Now().AddDate(0, 0, -30)
	params := map[string]interface{}{
		"created_at_min": after.Format(time.RFC3339),
	}
	return uc.SyncOrdersWithParams(ctx, integrationID, params)
}

func (uc *meliUseCase) SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	sellerID, err := extractSellerID(integration)
	if err != nil {
		return err
	}

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("ensuring valid token: %w", err)
	}

	queryParams := buildMeliQueryParams(params)

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Int64("seller_id", sellerID).
		Msg("Starting MercadoLibre order sync")

	go uc.syncOrdersAsync(context.Background(), integration, accessToken, sellerID, queryParams)

	return nil
}

func (uc *meliUseCase) syncOrdersAsync(ctx context.Context, integration *domain.Integration, accessToken string, sellerID int64, params *domain.GetOrdersParams) {
	start := time.Now()
	uc.emitOrderSyncEvent(ctx, integration, "integration.sync.started", map[string]interface{}{})

	totalSynced := 0
	offset := 0
	if params.Limit == 0 {
		params.Limit = 50
	}

	integrationID := fmt.Sprintf("%d", integration.ID)

	for {
		params.Offset = offset

		result, rawOrders, err := uc.client.GetOrders(ctx, accessToken, sellerID, params)
		if err != nil {
			if err == domain.ErrTokenExpired {
				uc.logger.Info(ctx).Msg("Token expired during sync, refreshing...")
				newToken, refreshErr := uc.EnsureValidToken(ctx, integrationID)
				if refreshErr != nil {
					uc.logger.Error(ctx).Err(refreshErr).Msg("Token refresh failed during sync")
					break
				}
				accessToken = newToken
				continue
			}

			uc.logger.Error(ctx).Err(err).
				Int("offset", offset).
				Msg("Error fetching MeLi orders page")
			uc.emitOrderSyncEvent(ctx, integration, "integration.sync.failed", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if len(result.Orders) == 0 {
			break
		}

		for i := range result.Orders {
			order := result.Orders[i]
			var rawJSON []byte
			if i < len(rawOrders) {
				rawJSON = rawOrders[i]
			}

			var shippingDetail *domain.MeliShippingDetail
			if order.Shipping != nil && order.Shipping.ID > 0 {
				detail, shErr := uc.client.GetShipmentDetail(ctx, accessToken, order.Shipping.ID)
				if shErr != nil {
					uc.logger.Warn(ctx).Err(shErr).
						Int64("shipment_id", order.Shipping.ID).
						Msg("Failed to fetch shipment detail during sync")
				} else {
					shippingDetail = detail
				}
			}

			uc.enrichBillingInfo(ctx, accessToken, &order)

			dto := mapper.MapMeliOrderToProbability(&order, shippingDetail, rawJSON)
			dto.IntegrationID = integration.ID
			dto.BusinessID = integration.BusinessID

			if err := uc.publisher.Publish(ctx, dto); err != nil {
				uc.logger.Error(ctx).Err(err).
					Int64("order_id", order.ID).
					Msg("Error publishing MeLi order")
				continue
			}
			totalSynced++
		}

		uc.logger.Info(ctx).
			Int("offset", offset).
			Int("orders_in_page", len(result.Orders)).
			Int("total", result.Total).
			Msg("MeLi orders page synced")

		if offset+params.Limit >= result.Total {
			break
		}

		offset += params.Limit
		time.Sleep(1 * time.Second)
	}

	uc.logger.Info(ctx).
		Int("total_synced", totalSynced).
		Uint("integration_id", integration.ID).
		Msg("MercadoLibre order sync completed")

	uc.emitOrderSyncEvent(ctx, integration, "integration.sync.completed", map[string]interface{}{
		"total_fetched": totalSynced,
		"duration":      time.Since(start).Round(time.Millisecond).String(),
	})
}

func extractSellerID(integration *domain.Integration) (int64, error) {
	if v, ok := integration.Config["seller_id"]; ok {
		switch id := v.(type) {
		case float64:
			return int64(id), nil
		case int64:
			return id, nil
		case string:
			if id != "" {
				var parsed int64
				_, err := fmt.Sscanf(id, "%d", &parsed)
				if err == nil {
					return parsed, nil
				}
			}
		}
	}

	if integration.StoreID != "" {
		var parsed int64
		_, err := fmt.Sscanf(integration.StoreID, "%d", &parsed)
		if err == nil {
			return parsed, nil
		}
	}

	return 0, domain.ErrSellerIDNotFound
}

func buildMeliQueryParams(params interface{}) *domain.GetOrdersParams {
	qp := &domain.GetOrdersParams{
		Sort:  "date_desc",
		Limit: 50,
	}

	m, ok := params.(map[string]interface{})
	if !ok {
		return qp
	}

	if v, ok := m["created_at_min"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			qp.DateFrom = &t
		}
	}
	if v, ok := m["created_at_max"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			qp.DateTo = &t
		}
	}
	if v, ok := m["status"].(string); ok && v != "" {
		qp.Status = v
	}

	return qp
}
