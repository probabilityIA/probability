package usecases

import (
	"context"
	"fmt"
)

// DeleteWebhook elimina un webhook de Shopify y actualiza el config de la integración
func (uc *SyncOrdersUseCase) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
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

	// Eliminar webhook usando el cliente
	if err := uc.shopifyClient.DeleteWebhook(ctx, storeName, accessToken, webhookID); err != nil {
		return fmt.Errorf("error al eliminar webhook: %w", err)
	}

	// Actualizar el config removiendo el webhook ID de la lista
	webhookIDs, ok := config["webhook_ids"].([]interface{})
	if ok {
		// Filtrar el webhookID eliminado
		newWebhookIDs := []interface{}{}
		for _, id := range webhookIDs {
			idStr := fmt.Sprintf("%v", id)
			if idStr != webhookID {
				newWebhookIDs = append(newWebhookIDs, id)
			}
		}

		// Si no quedan webhooks, limpiar la configuración
		if len(newWebhookIDs) == 0 {
			configUpdate := map[string]interface{}{
				"webhook_ids":        []string{},
				"webhook_configured": false,
			}
			if err := uc.integrationService.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
				return fmt.Errorf("error al actualizar config: %w", err)
			}
		} else {
			configUpdate := map[string]interface{}{
				"webhook_ids": newWebhookIDs,
			}
			if err := uc.integrationService.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
				return fmt.Errorf("error al actualizar config: %w", err)
			}
		}
	}

	return nil
}













