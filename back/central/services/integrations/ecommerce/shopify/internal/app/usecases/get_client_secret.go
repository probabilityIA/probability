package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

// GetClientSecretByShopDomain recupera el Client Secret de una integración de Shopify por su dominio
func (uc *SyncOrdersUseCase) GetClientSecretByShopDomain(ctx context.Context, shopDomain string) (string, error) {
	// Buscar la integración por store_id (que es el shop domain)
	integration, err := uc.integrationService.GetIntegrationByExternalID(ctx, shopDomain, domain.IntegrationTypeID)
	if err != nil {
		return "", fmt.Errorf("error al buscar integración para el dominio %s: %w", shopDomain, err)
	}

	if integration == nil {
		return "", fmt.Errorf("integración no encontrada para el dominio: %s", shopDomain)
	}

	// Recuperar el secreto usando el servicio de descifrado
	secret, err := uc.integrationService.DecryptCredential(ctx, fmt.Sprintf("%d", integration.ID), "client_secret")
	if err != nil {
		return "", fmt.Errorf("error al descifrar client_secret para %s: %w", shopDomain, err)
	}

	return secret, nil
}
