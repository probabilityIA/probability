package usecases

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

// CreateWebhook crea un webhook en Shopify para la integración y actualiza el config con la información
// Primero verifica si existen webhooks con la misma URL y los elimina antes de crear nuevos
func (uc *SyncOrdersUseCase) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (*domain.CreateWebhookResult, error) {
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

	// Construir nuestra URL del webhook
	webhookURL := fmt.Sprintf("%s/integrations/shopify/webhook", baseURL)

	// Validar si la URL es localhost (entorno de pruebas)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error al parsear baseURL: %w", err)
	}

	hostname := strings.ToLower(parsedURL.Hostname())
	if strings.Contains(hostname, "localhost") || strings.Contains(hostname, "127.0.0.1") || strings.Contains(hostname, "::1") {
		return &domain.CreateWebhookResult{
			ExistingWebhooks: []domain.WebhookInfo{},
			DeletedWebhooks:  []domain.WebhookInfo{},
			CreatedWebhooks:  []string{},
			WebhookURL:       webhookURL,
		}, fmt.Errorf("no se pueden crear webhooks en entorno de pruebas (localhost). La URL del webhook sería: %s", webhookURL)
	}

	// Verificar si existen webhooks con nuestra URL
	existingWebhooks, err := uc.VerifyWebhooksByURL(ctx, integrationID, baseURL)
	if err != nil {
		return nil, fmt.Errorf("error al verificar webhooks existentes: %w", err)
	}

	result := &domain.CreateWebhookResult{
		ExistingWebhooks: existingWebhooks,
		DeletedWebhooks:  make([]domain.WebhookInfo, 0),
		CreatedWebhooks:  make([]string, 0),
		WebhookURL:       webhookURL,
	}

	// Eliminar solo los webhooks que coinciden con nuestra URL
	for _, webhook := range existingWebhooks {
		if err := uc.shopifyClient.DeleteWebhook(ctx, storeName, accessToken, webhook.ID); err != nil {
			// Log del error pero continuamos con los demás
			// No fallamos si no podemos eliminar un webhook existente
			continue
		}
		result.DeletedWebhooks = append(result.DeletedWebhooks, webhook)
	}

	// Eventos que necesitamos registrar
	events := []string{
		"orders/create",
		"orders/updated",
		"orders/paid",
		"orders/cancelled",
		"orders/fulfilled",
		"orders/partially_fulfilled",
	}

	// Crear webhooks para todos los eventos
	webhookConfigured := true
	for _, event := range events {
		webhookID, err := uc.shopifyClient.CreateWebhook(ctx, storeName, accessToken, webhookURL, event)
		if err != nil {
			// Si falla la creación, marcamos como no configurado pero continuamos con los demás
			webhookConfigured = false
			// Log del error pero continuamos
			continue
		}
		result.CreatedWebhooks = append(result.CreatedWebhooks, webhookID)
	}

	// Si no se creó ningún webhook, retornar error
	if len(result.CreatedWebhooks) == 0 {
		return result, fmt.Errorf("no se pudo crear ningún webhook en Shopify")
	}

	// Actualizar el config con la información del webhook
	configUpdate := map[string]interface{}{
		"webhook_url":        webhookURL,
		"webhook_configured": webhookConfigured,
		"webhook_ids":        result.CreatedWebhooks,
	}

	// Hacer merge con el config existente
	if err := uc.integrationService.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
		return result, fmt.Errorf("error al actualizar config de la integración: %w", err)
	}

	return result, nil
}
