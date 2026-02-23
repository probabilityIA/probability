package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

// ListWebhooks lista todos los webhooks de una integración de Shopify
func (uc *SyncOrdersUseCase) ListWebhooks(ctx context.Context, integrationID string) ([]domain.WebhookInfo, error) {
	// Obtener la integración
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener las credenciales
	accessToken, err := uc.integrationService.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return nil, fmt.Errorf("error al obtener access_token: %w", err)
	}

	// Obtener el store_name del config
	storeName, ok := integration.Config["store_name"].(string)
	if !ok || storeName == "" {
		return nil, fmt.Errorf("store_name no encontrado en la configuración")
	}

	// Listar webhooks usando el cliente
	webhooks, err := uc.shopifyClient.ListWebhooks(ctx, storeName, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error al listar webhooks: %w", err)
	}

	return webhooks, nil
}













