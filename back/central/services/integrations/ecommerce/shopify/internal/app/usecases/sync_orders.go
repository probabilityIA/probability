package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	integrationevents "github.com/secamc93/probability/back/central/services/integrations/events"
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

	fmt.Printf("[SyncOrders] Starting sync for integration %s. Params: CreatedAtMin=%v, CreatedAtMax=%v, Status=%s, FinancialStatus=%s, FulfillmentStatus=%s\n",
		integrationID, params.CreatedAtMin, params.CreatedAtMax, params.Status, params.FinancialStatus, params.FulfillmentStatus)

	// #region agent log
	if f, err := os.OpenFile("/home/cam/Desktop/probability/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		logData, _ := json.Marshal(map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "A",
			"location":     "sync_orders.go:76",
			"message":      "SyncOrdersWithParams - Date params being sent to Shopify API",
			"data": map[string]interface{}{
				"integration_id":     integrationID,
				"created_at_min":     params.CreatedAtMin,
				"created_at_max":     params.CreatedAtMax,
				"status":             params.Status,
				"financial_status":   params.FinancialStatus,
				"fulfillment_status": params.FulfillmentStatus,
				"store_domain":       storeDomain,
			},
			"timestamp": time.Now().UnixMilli(),
		})
		f.WriteString(string(logData) + "\n")
		f.Close()
	}
	// #endregion

	// Publicar evento de inicio de sincronización
	integrationIDUint, _ := strconv.ParseUint(integrationID, 10, 32)
	integrationevents.PublishSyncStarted(
		ctx,
		uint(integrationIDUint),
		integration.BusinessID,
		"shopify",
		syncParams.CreatedAtMin,
		syncParams.CreatedAtMax,
		syncParams.Status,
		syncParams.FinancialStatus,
		syncParams.FulfillmentStatus,
	)

	go func() {
		ctx := context.Background()
		startTime := time.Now()
		var totalOrders, createdOrders, updatedOrders, rejectedOrders int

		if err := uc.GetOrders(ctx, integration, storeDomain, accessToken, params); err != nil {
			fmt.Printf("[SyncOrders] Error in GetOrders: %v\n", err)
			// Publicar evento de sincronización fallida
			integrationevents.PublishSyncFailed(
				ctx,
				uint(integrationIDUint),
				integration.BusinessID,
				"shopify",
				err.Error(),
			)
			return
		}

		// Publicar evento de sincronización completada
		// Nota: Los contadores (totalOrders, createdOrders, etc.) se actualizarán cuando se procesen las órdenes
		// Por ahora, publicamos con valores iniciales. Se puede mejorar para rastrear estos valores.
		integrationevents.PublishSyncCompleted(
			ctx,
			uint(integrationIDUint),
			integration.BusinessID,
			"shopify",
			totalOrders,
			createdOrders,
			updatedOrders,
			rejectedOrders,
			time.Since(startTime),
		)
	}()

	return nil
}
