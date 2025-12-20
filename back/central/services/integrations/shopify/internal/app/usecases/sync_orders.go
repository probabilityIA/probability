package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases/utils"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
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

	go func() {
		ctx := context.Background()
		if err := uc.GetOrders(ctx, integration, storeDomain, accessToken, params); err != nil {
			fmt.Printf("[SyncOrders] Error in GetOrders: %v\n", err)
		}
	}()

	return nil
}
