package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// SyncOrders sincroniza órdenes de WooCommerce de los últimos 30 días.
func (uc *wooCommerceUseCase) SyncOrders(ctx context.Context, integrationID string) error {
	after := time.Now().AddDate(0, 0, -30)
	params := map[string]interface{}{
		"created_at_min": after.Format(time.RFC3339),
	}
	return uc.SyncOrdersWithParams(ctx, integrationID, params)
}

// SyncOrdersWithParams sincroniza órdenes con parámetros personalizados.
// Soporta: created_at_min, created_at_max, status (string).
func (uc *wooCommerceUseCase) SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error {
	// 1. Obtener integración
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	// 2. Extraer configuración
	storeURL, err := extractString(integration.Config, "store_url")
	if err != nil {
		return domain.ErrMissingStoreURL
	}

	// 3. Descifrar credenciales
	consumerKey, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_key")
	if err != nil {
		return fmt.Errorf("decrypting consumer_key: %w", err)
	}
	consumerSecret, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_secret")
	if err != nil {
		return fmt.Errorf("decrypting consumer_secret: %w", err)
	}

	// 4. Construir parámetros de consulta
	queryParams := buildQueryParams(params)

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("store_url", storeURL).
		Msg("Starting WooCommerce order sync")

	// 5. Lanzar sincronización en background
	go uc.syncOrdersAsync(context.Background(), integration, storeURL, consumerKey, consumerSecret, queryParams)

	return nil
}

func (uc *wooCommerceUseCase) syncOrdersAsync(ctx context.Context, integration *domain.Integration, storeURL, consumerKey, consumerSecret string, params *domain.GetOrdersParams) {
	totalSynced := 0
	page := 1
	if params.PerPage == 0 {
		params.PerPage = 100
	}

	for {
		params.Page = page

		result, rawOrders, err := uc.client.GetOrders(ctx, storeURL, consumerKey, consumerSecret, params)
		if err != nil {
			uc.logger.Error(ctx).Err(err).
				Int("page", page).
				Msg("Error fetching WooCommerce orders page")
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

			dto := mapper.MapWooOrderToProbability(&order, rawJSON)
			dto.IntegrationID = integration.ID
			dto.BusinessID = integration.BusinessID

			if err := uc.publisher.Publish(ctx, dto); err != nil {
				uc.logger.Error(ctx).Err(err).
					Str("order_number", order.Number).
					Msg("Error publishing WooCommerce order")
				continue
			}
			totalSynced++
		}

		uc.logger.Info(ctx).
			Int("page", page).
			Int("orders_in_page", len(result.Orders)).
			Int("total_pages", result.TotalPages).
			Msg("WooCommerce orders page synced")

		if page >= result.TotalPages {
			break
		}

		page++
		// Rate limiting: 500ms entre páginas
		time.Sleep(500 * time.Millisecond)
	}

	uc.logger.Info(ctx).
		Int("total_synced", totalSynced).
		Uint("integration_id", integration.ID).
		Msg("WooCommerce order sync completed")
}

// buildQueryParams construye los parámetros de consulta a partir de un mapa genérico.
func buildQueryParams(params interface{}) *domain.GetOrdersParams {
	qp := &domain.GetOrdersParams{
		OrderBy: "date",
		Order:   "desc",
	}

	m, ok := params.(map[string]interface{})
	if !ok {
		return qp
	}

	if v, ok := m["created_at_min"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			qp.After = &t
		}
	}
	if v, ok := m["created_at_max"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			qp.Before = &t
		}
	}
	if v, ok := m["status"].(string); ok && v != "" {
		qp.Status = v
	}

	return qp
}
