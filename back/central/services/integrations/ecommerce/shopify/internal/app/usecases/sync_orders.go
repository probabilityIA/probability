package usecases

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/utils"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

// SyncOrders sincroniza órdenes con parámetros por defecto (últimos 30 días)
func (uc *SyncOrdersUseCase) SyncOrders(ctx context.Context, integrationID string) error {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	params := &domain.SyncOrdersParams{
		CreatedAtMin:      &thirtyDaysAgo,
		Status:            domain.OrderStatusAny,
		FinancialStatus:   domain.FinancialStatusAny,
		FulfillmentStatus: domain.FulfillmentStatusAny,
	}
	return uc.SyncOrdersWithParams(ctx, integrationID, params)
}

// SyncOrdersWithParams sincroniza órdenes con parámetros personalizados
func (uc *SyncOrdersUseCase) SyncOrdersWithParams(ctx context.Context, integrationID string, syncParams *domain.SyncOrdersParams) error {
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	config, err := utils.NormalizeConfig(integration.Config, integration.Name)
	if err != nil {
		return err
	}

	storeDomain, err := utils.ExtractStoreName(config, integration.Name)
	if err != nil {
		return fmt.Errorf("failed to extract store name: %w", err)
	}

	// En modo test, usar la URL de pruebas (base_url_test) en vez del dominio de Shopify
	storeDomain = utils.ResolveEffectiveStoreDomain(integration, storeDomain)

	accessToken, err := utils.GetAccessToken(ctx, uc.integrationService, integrationID)
	if err != nil {
		return err
	}

	// Construir parámetros para la API de Shopify
	params := &domain.GetOrdersParams{
		Limit: 250,
	}

	// Aplicar filtros de fecha
	if syncParams.CreatedAtMin != nil {
		params.CreatedAtMin = syncParams.CreatedAtMin
	}
	if syncParams.CreatedAtMax != nil {
		params.CreatedAtMax = syncParams.CreatedAtMax
	}

	// Aplicar filtros de estado
	if syncParams.Status != "" {
		params.Status = syncParams.Status
	} else {
		params.Status = domain.OrderStatusAny
	}

	if syncParams.FinancialStatus != "" {
		params.FinancialStatus = syncParams.FinancialStatus
	} else {
		params.FinancialStatus = domain.FinancialStatusAny
	}

	if syncParams.FulfillmentStatus != "" {
		params.FulfillmentStatus = syncParams.FulfillmentStatus
	}

	// Parsear integration ID a uint para eventos SSE
	integrationIDUint, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("integration_id inválido")
		return fmt.Errorf("%w: %s", domain.ErrInvalidIntegrationID, integrationID)
	}
	intIDUint := uint(integrationIDUint)

	uc.log.Info(ctx).
		Str("integration_id", integrationID).
		Interface("created_at_min", params.CreatedAtMin).
		Interface("created_at_max", params.CreatedAtMax).
		Str("status", params.Status).
		Str("financial_status", params.FinancialStatus).
		Str("fulfillment_status", params.FulfillmentStatus).
		Msg("Iniciando sincronización de órdenes")

	// Publicar evento de inicio de sincronización
	uc.syncEventPublisher.PublishSyncEvent(ctx, intIDUint, integration.BusinessID, "integration.sync.started", map[string]interface{}{
		"integration_id":   intIDUint,
		"integration_type": "shopify",
		"params": map[string]interface{}{
			"created_at_min":     syncParams.CreatedAtMin,
			"created_at_max":     syncParams.CreatedAtMax,
			"status":             syncParams.Status,
			"financial_status":   syncParams.FinancialStatus,
			"fulfillment_status": syncParams.FulfillmentStatus,
		},
		"started_at": time.Now(),
	})

	startTime := time.Now()

	totalFetched, err := uc.GetOrders(ctx, integration, storeDomain, accessToken, params)
	if err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Error en GetOrders")
		// Publicar evento de sincronización fallida
		uc.syncEventPublisher.PublishSyncEvent(ctx, intIDUint, integration.BusinessID, "integration.sync.failed", map[string]interface{}{
			"integration_id":   intIDUint,
			"integration_type": "shopify",
			"error":            err.Error(),
			"failed_at":        time.Now(),
		})
		return err
	}

	// Publicar evento de sincronización completada
	// total_fetched = órdenes publicadas a la cola (aún no procesadas por consumers)
	uc.syncEventPublisher.PublishSyncEvent(ctx, intIDUint, integration.BusinessID, "integration.sync.completed", map[string]interface{}{
		"integration_id":   intIDUint,
		"integration_type": "shopify",
		"total_fetched":    totalFetched,
		"duration":         time.Since(startTime).String(),
		"completed_at":     time.Now(),
	})

	return nil
}
