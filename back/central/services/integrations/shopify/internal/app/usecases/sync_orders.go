package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases/utils"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) SyncOrders(ctx context.Context, integrationID string) error {
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

	fifteenDaysAgo := time.Now().AddDate(0, 0, -15)
	params := &domain.GetOrdersParams{
		Status:          "any",
		Limit:           250,
		CreatedAtMin:    &fifteenDaysAgo,
		FinancialStatus: "paid",
	}

	fmt.Printf("[SyncOrders] Starting sync for integration %s. Params: CreatedAtMin=%v, Status=%s, Limit=%d\n",
		integrationID, params.CreatedAtMin, params.Status, params.Limit)

	go func() {
		ctx := context.Background()
		if err := uc.GetOrders(ctx, integration, storeDomain, accessToken, params); err != nil {
			fmt.Printf("[SyncOrders] Error in GetOrders: %v\n", err)
		}
	}()

	return nil
}
