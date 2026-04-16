package usecases

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

// SyncOrders sincroniza órdenes de VTEX de los últimos 30 días.
func (uc *vtexUseCase) SyncOrders(ctx context.Context, integrationID string) error {
	after := time.Now().AddDate(0, 0, -30)
	params := map[string]interface{}{
		"created_at_min": after.Format(time.RFC3339),
	}
	return uc.SyncOrdersWithParams(ctx, integrationID, params)
}

// SyncOrdersWithParams sincroniza órdenes con parámetros personalizados.
// Soporta: created_at_min, created_at_max, status (string).
func (uc *vtexUseCase) SyncOrdersWithParams(ctx context.Context, integrationID string, params interface{}) error {
	// 1. Obtener integración
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	// 2. Obtener credenciales
	storeURL, apiKey, apiToken, err := uc.getCredentials(ctx, integration, integrationID)
	if err != nil {
		return err
	}

	// 3. Construir filtros
	filters := buildVTEXFilters(params)

	uc.logger.Info(ctx).
		Str("integration_id", integrationID).
		Str("store_url", storeURL).
		Msg("Starting VTEX order sync")

	// 4. Lanzar sincronización en background
	go uc.syncOrdersAsync(context.Background(), integration, storeURL, apiKey, apiToken, filters)

	return nil
}

func (uc *vtexUseCase) syncOrdersAsync(ctx context.Context, integration *domain.Integration, storeURL, apiKey, apiToken string, filters map[string]string) {
	totalSynced := 0
	page := 1
	perPage := 15 // VTEX default

	for {
		result, err := uc.client.GetOrders(ctx, storeURL, apiKey, apiToken, page, perPage, filters)
		if err != nil {
			uc.logger.Error(ctx).Err(err).
				Int("page", page).
				Msg("Error fetching VTEX orders page")
			break
		}

		if len(result.List) == 0 {
			break
		}

		for _, summary := range result.List {
			// Obtener detalle completo de cada orden
			order, rawJSON, err := uc.client.GetOrderByID(ctx, storeURL, apiKey, apiToken, summary.OrderID)
			if err != nil {
				uc.logger.Error(ctx).Err(err).
					Str("order_id", summary.OrderID).
					Msg("Error fetching VTEX order detail")
				continue
			}

			dto := mapper.MapVTEXOrderToProbability(order, rawJSON)
			dto.IntegrationID = integration.ID
			dto.BusinessID = integration.BusinessID

			if err := uc.publisher.Publish(ctx, dto); err != nil {
				uc.logger.Error(ctx).Err(err).
					Str("order_id", summary.OrderID).
					Msg("Error publishing VTEX order")
				continue
			}
			totalSynced++
		}

		uc.logger.Info(ctx).
			Int("page", page).
			Int("orders_in_page", len(result.List)).
			Int("total", result.Paging.Total).
			Msg("VTEX orders page synced")

		// Si ya estamos en la última página, terminar
		if page >= result.Paging.Pages {
			break
		}

		page++
		// Rate limiting: 500ms entre páginas
		time.Sleep(500 * time.Millisecond)
	}

	uc.logger.Info(ctx).
		Int("total_synced", totalSynced).
		Uint("integration_id", integration.ID).
		Msg("VTEX order sync completed")
}

// getCredentials obtiene store_url, api_key y api_token de una integración.
func (uc *vtexUseCase) getCredentials(ctx context.Context, integration *domain.Integration, integrationID string) (storeURL, apiKey, apiToken string, err error) {
	storeURL, err = extractString(integration.Config, "store_url")
	if err != nil {
		return "", "", "", domain.ErrMissingStoreURL
	}

	apiKey, err = uc.service.DecryptCredential(ctx, integrationID, "api_key")
	if err != nil {
		return "", "", "", fmt.Errorf("decrypting api_key: %w", err)
	}
	if apiKey == "" {
		return "", "", "", domain.ErrMissingAPIKey
	}

	apiToken, err = uc.service.DecryptCredential(ctx, integrationID, "api_token")
	if err != nil {
		return "", "", "", fmt.Errorf("decrypting api_token: %w", err)
	}
	if apiToken == "" {
		return "", "", "", domain.ErrMissingAPIToken
	}

	return storeURL, apiKey, apiToken, nil
}

// buildVTEXFilters construye los filtros de la API de VTEX a partir de un mapa genérico.
func buildVTEXFilters(params interface{}) map[string]string {
	filters := make(map[string]string)

	m, ok := params.(map[string]interface{})
	if !ok {
		return filters
	}

	// Convertir created_at_min/max a f_creationDate de VTEX
	var dateFrom, dateTo string
	if v, ok := m["created_at_min"].(string); ok && v != "" {
		dateFrom = v
	}
	if v, ok := m["created_at_max"].(string); ok && v != "" {
		dateTo = v
	} else {
		dateTo = time.Now().Format(time.RFC3339)
	}

	if dateFrom != "" {
		// VTEX format: f_creationDate=creationDate:[2026-01-01T00:00:00.000Z TO 2026-02-24T23:59:59.999Z]
		filters["f_creationDate"] = url.QueryEscape(fmt.Sprintf("creationDate:[%s TO %s]", dateFrom, dateTo))
	}

	if v, ok := m["status"].(string); ok && v != "" {
		filters["f_status"] = v
	}

	return filters
}
