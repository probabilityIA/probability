package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

// SyncOrders sincroniza órdenes de MercadoLibre de los últimos 30 días.
func (uc *meliUseCase) SyncOrders(ctx context.Context, integrationID string) error {
	after := time.Now().AddDate(0, 0, -30)
	params := map[string]interface{}{
		"created_at_min": after.Format(time.RFC3339),
	}
	return uc.SyncOrdersWithParams(ctx, integrationID, params)
}

// SyncOrdersWithParams sincroniza órdenes con parámetros personalizados.
// Soporta: created_at_min, created_at_max, status (string).
func (uc *meliUseCase) SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error {
	// 1. Obtener integración
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	// 2. Obtener seller_id del config o store_id
	sellerID, err := extractSellerID(integration)
	if err != nil {
		return err
	}

	// 3. Obtener access_token vigente
	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("ensuring valid token: %w", err)
	}

	// 4. Construir query params
	queryParams := buildMeliQueryParams(params)

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Int64("seller_id", sellerID).
		Msg("Starting MercadoLibre order sync")

	// 5. Lanzar sincronización en background
	go uc.syncOrdersAsync(context.Background(), integration, accessToken, sellerID, queryParams)

	return nil
}

func (uc *meliUseCase) syncOrdersAsync(ctx context.Context, integration *domain.Integration, accessToken string, sellerID int64, params *domain.GetOrdersParams) {
	totalSynced := 0
	offset := 0
	if params.Limit == 0 {
		params.Limit = 50 // MeLi max per page
	}

	integrationID := fmt.Sprintf("%d", integration.ID)

	for {
		params.Offset = offset

		result, rawOrders, err := uc.client.GetOrders(ctx, accessToken, sellerID, params)
		if err != nil {
			// Si el token expiró durante la sync, intentar refrescar
			if err == domain.ErrTokenExpired {
				uc.logger.Info(ctx).Msg("Token expired during sync, refreshing...")
				newToken, refreshErr := uc.EnsureValidToken(ctx, integrationID)
				if refreshErr != nil {
					uc.logger.Error(ctx).Err(refreshErr).Msg("Token refresh failed during sync")
					break
				}
				accessToken = newToken
				continue // Reintentar la misma página
			}

			uc.logger.Error(ctx).Err(err).
				Int("offset", offset).
				Msg("Error fetching MeLi orders page")
			break
		}

		if len(result.Orders) == 0 {
			break
		}

		for i, order := range result.Orders {
			var rawJSON []byte
			if i < len(rawOrders) {
				rawJSON = rawOrders[i]
			}

			// Obtener detalles de envío si hay shipping
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

		// MeLi paginación: si offset + limit >= total, terminamos
		if offset+params.Limit >= result.Total {
			break
		}

		offset += params.Limit
		// Rate limiting: 1s entre páginas (MeLi es más agresivo que WooCommerce)
		time.Sleep(1 * time.Second)
	}

	uc.logger.Info(ctx).
		Int("total_synced", totalSynced).
		Uint("integration_id", integration.ID).
		Msg("MercadoLibre order sync completed")
}

// extractSellerID obtiene el seller_id de una integración.
// Busca primero en config["seller_id"], luego usa StoreID.
func extractSellerID(integration *domain.Integration) (int64, error) {
	// Opción 1: config["seller_id"]
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

	// Opción 2: StoreID
	if integration.StoreID != "" {
		var parsed int64
		_, err := fmt.Sscanf(integration.StoreID, "%d", &parsed)
		if err == nil {
			return parsed, nil
		}
	}

	return 0, domain.ErrSellerIDNotFound
}

// buildMeliQueryParams construye los parámetros de consulta a partir de un mapa genérico.
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
