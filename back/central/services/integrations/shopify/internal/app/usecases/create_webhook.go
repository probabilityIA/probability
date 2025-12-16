package usecases

import (
	"context"
	"fmt"
)

// CreateWebhook crea un webhook en Shopify para la integración y actualiza el config con la información
func (uc *SyncOrdersUseCase) CreateWebhook(ctx context.Context, integrationID string, baseURL string) error {
	// Obtener la integración
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener las credenciales
	accessToken, err := uc.integrationService.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return fmt.Errorf("error al obtener access_token: %w", err)
	}

	// Obtener el store_name del config
	config, ok := integration.Config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("config no es un map válido")
	}

	storeName, ok := config["store_name"].(string)
	if !ok || storeName == "" {
		return fmt.Errorf("store_name no encontrado en la configuración")
	}

	// Obtener la URL del webhook usando GetWebhookURL del core
	// Para esto necesitamos acceder al coreIntegration, así que lo haremos desde el caso de uso
	// Por ahora, construiremos la URL directamente basándonos en la lógica existente
	webhookURL := fmt.Sprintf("%s/integrations/shopify/webhook", baseURL)

	// Eventos que necesitamos registrar
	events := []string{
		"orders/create",
		"orders/updated",
		"orders/paid",
		"orders/cancelled",
		"orders/fulfilled",
		"orders/partially_fulfilled",
	}

	// Intentar crear webhooks para todos los eventos
	// Shopify permite múltiples webhooks para la misma URL con diferentes eventos
	webhookConfigured := true
	webhookIDs := make([]string, 0)

	for _, event := range events {
		webhookID, err := uc.shopifyClient.CreateWebhook(ctx, storeName, accessToken, webhookURL, event)
		if err != nil {
			// Si falla la creación, marcamos como no configurado pero continuamos con los demás
			webhookConfigured = false
			// Log del error pero continuamos
			continue
		}
		webhookIDs = append(webhookIDs, webhookID)
	}

	// Si no se creó ningún webhook, retornar error
	if len(webhookIDs) == 0 {
		return fmt.Errorf("no se pudo crear ningún webhook en Shopify")
	}

	// Actualizar el config con la información del webhook
	configUpdate := map[string]interface{}{
		"webhook_url":         webhookURL,
		"webhook_configured":  webhookConfigured,
		"webhook_ids":         webhookIDs,
	}

	// Hacer merge con el config existente
	if err := uc.integrationService.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
		return fmt.Errorf("error al actualizar config de la integración: %w", err)
	}

	return nil
}

