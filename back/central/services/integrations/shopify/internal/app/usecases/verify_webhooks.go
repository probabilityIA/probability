package usecases

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

// VerifyWebhooksByURL busca webhooks existentes que coincidan con nuestra URL
// Retorna los webhooks encontrados que coinciden con nuestra URL generada
func (uc *SyncOrdersUseCase) VerifyWebhooksByURL(ctx context.Context, integrationID string, baseURL string) ([]domain.WebhookInfo, error) {
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
	config, ok := integration.Config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("config no es un map válido")
	}

	storeName, ok := config["store_name"].(string)
	if !ok || storeName == "" {
		return nil, fmt.Errorf("store_name no encontrado en la configuración")
	}

	// Construir nuestra URL de webhook
	ourWebhookURL := fmt.Sprintf("%s/integrations/shopify/webhook", baseURL)

	// Normalizar URLs para comparación (sin trailing slash, en minúsculas)
	ourURLParsed, err := url.Parse(ourWebhookURL)
	if err != nil {
		return nil, fmt.Errorf("error al parsear nuestra URL de webhook: %w", err)
	}
	// Normalizar: remover trailing slash, convertir a minúsculas
	ourURLNormalized := strings.TrimSuffix(strings.ToLower(ourURLParsed.String()), "/")

	// Listar todos los webhooks de Shopify
	allWebhooks, err := uc.shopifyClient.ListWebhooks(ctx, storeName, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error al listar webhooks: %w", err)
	}

	// Filtrar solo los webhooks que coinciden con nuestra URL
	matchingWebhooks := make([]domain.WebhookInfo, 0)
	for _, webhook := range allWebhooks {
		webhookURLParsed, err := url.Parse(webhook.Address)
		if err != nil {
			// Si no podemos parsear la URL, la ignoramos
			continue
		}
		// Normalizar: remover trailing slash, convertir a minúsculas
		webhookURLNormalized := strings.TrimSuffix(strings.ToLower(webhookURLParsed.String()), "/")

		// Comparar URLs normalizadas
		if webhookURLNormalized == ourURLNormalized {
			matchingWebhooks = append(matchingWebhooks, webhook)
		}
	}

	return matchingWebhooks, nil
}
