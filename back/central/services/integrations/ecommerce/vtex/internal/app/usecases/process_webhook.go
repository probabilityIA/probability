package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

// ProcessWebhook procesa un webhook de VTEX (cambio de estado de orden).
// Recibe el payload parseado, obtiene la orden completa de la API y la publica a la cola.
func (uc *vtexUseCase) ProcessWebhook(ctx context.Context, payload *domain.VTEXWebhookPayload) error {
	if payload.OrderID == "" {
		uc.logger.Warn(ctx).Msg("Ignoring VTEX webhook with empty OrderId")
		return nil
	}

	uc.logger.Info(ctx).
		Str("order_id", payload.OrderID).
		Str("state", payload.State).
		Str("domain", payload.Domain).
		Msg("Processing VTEX webhook")

	// 1. Buscar integración por account name (del webhook Origin)
	accountName := ""
	if payload.Origin != nil {
		accountName = payload.Origin.Account
	}

	// Buscar la integración. Necesitamos iterar posibles integraciones.
	// El webhook no trae integrationID, así que buscamos por store_url que contiene el account name.
	// Si tenemos el account name, intentamos encontrar la integración que lo tenga en su store_url.
	// Fallback: buscar por el orderId prefix (que contiene el account en VTEX format).
	integration, integrationID, err := uc.findIntegrationForWebhook(ctx, accountName, payload.OrderID)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("account", accountName).
			Str("order_id", payload.OrderID).
			Msg("Failed to find integration for VTEX webhook")
		return fmt.Errorf("finding integration: %w", err)
	}
	if integration == nil {
		uc.logger.Warn(ctx).
			Str("account", accountName).
			Msg("No integration found for VTEX webhook")
		return domain.ErrIntegrationNotFound
	}

	// 2. Obtener credenciales
	storeURL, apiKey, apiToken, err := uc.getCredentials(ctx, integration, integrationID)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Uint("integration_id", integration.ID).
			Msg("Failed to get VTEX credentials")
		return fmt.Errorf("getting credentials: %w", err)
	}

	// 3. Obtener orden completa de la API
	order, rawJSON, err := uc.client.GetOrderByID(ctx, storeURL, apiKey, apiToken, payload.OrderID)
	if err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("order_id", payload.OrderID).
			Msg("Failed to fetch VTEX order detail")
		return fmt.Errorf("fetching order: %w", err)
	}

	// 4. Mapear a DTO canónico
	dto := mapper.MapVTEXOrderToProbability(order, rawJSON)
	dto.IntegrationID = integration.ID
	dto.BusinessID = integration.BusinessID

	// 5. Publicar a la cola
	if err := uc.publisher.Publish(ctx, dto); err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("order_id", payload.OrderID).
			Msg("Failed to publish VTEX order to queue")
		return fmt.Errorf("publishing order: %w", err)
	}

	uc.logger.Info(ctx).
		Str("order_id", payload.OrderID).
		Str("status", order.Status).
		Uint("integration_id", integration.ID).
		Msg("VTEX order published successfully via webhook")

	return nil
}

// findIntegrationForWebhook busca la integración VTEX asociada al webhook.
// Intenta buscar por account name en la config store_url de las integraciones.
// Si no encuentra por account, intenta extraer el integrationID del orderId prefix.
func (uc *vtexUseCase) findIntegrationForWebhook(ctx context.Context, accountName string, orderID string) (*domain.Integration, string, error) {
	// VTEX webhooks incluyen el account name en Origin.Account.
	// Nuestra store_url tiene el formato: https://{accountName}.vtexcommercestable.com.br
	// Buscamos integraciones cuya store_url contenga el account name.

	// Para simplificar la búsqueda cuando hay múltiples integraciones VTEX,
	// el account name se almacena en config["account_name"] o se extrae de store_url.
	// Si solo hay una integración VTEX activa, la usamos directamente.

	// Por ahora, si tenemos account name, buscamos por store_url.
	// Si no, necesitaríamos un mecanismo diferente.

	// La implementación más práctica: guardar el account_name en config al crear la integración,
	// y buscar por ese campo. Como IIntegrationService solo nos permite buscar por ID,
	// una alternativa es almacenar el VTEX account como store_id de la integración.

	// Con la interfaz actual, necesitamos que el webhook incluya algún identificador que
	// podamos usar. En VTEX, el Origin.Account es el account name.
	// Si no está disponible, no podemos encontrar la integración.

	if accountName == "" {
		return nil, "", fmt.Errorf("vtex webhook missing Origin.Account, cannot identify integration")
	}

	// Buscar por account_name almacenado en config
	// Esto requiere que al crear la integración se guarde config["account_name"]
	// Por ahora delegamos al service que soporta búsqueda.

	// Workaround pragmático: el account_name está en store_id (se configura al crear la integración)
	// TODO: cuando IIntegrationService soporte búsqueda por store_id, usar ese método.
	// Por ahora, usamos un approach simplificado con config.

	// Usar el integrationID si está almacenado en config como webhook_integration_id
	// Esto es un placeholder hasta que tengamos búsqueda por store_id.

	// Nota: en la práctica, VTEX account_name se usa como store_id en la integración.
	// Para producción, necesitamos GetIntegrationByStoreID como tiene MeLi.
	// Por ahora, logueamos un warning y retornamos error.
	return nil, "", fmt.Errorf("VTEX webhook lookup by account_name '%s' not yet fully implemented: store the VTEX account name as store_id in the integration", accountName)
}
