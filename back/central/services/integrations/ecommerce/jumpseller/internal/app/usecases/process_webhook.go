package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases/mapper"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/client/response"
)

var ignoredWebhookEvents = map[string]bool{
	domain.EventProductStockUpdate: true,
	domain.EventOrderAbandoned:     true,
}

func (uc *jumpsellerUseCase) ResolveHooksToken(ctx context.Context, integrationID string) (string, error) {
	integration, cred, err := uc.resolveIntegration(ctx, integrationID)
	if err != nil {
		return "", err
	}

	if token, err := extractString(integration.Config, "hooks_token"); err == nil {
		return token, nil
	}

	storeInfo, err := uc.client.GetStoreInfo(ctx, cred)
	if err != nil {
		return "", fmt.Errorf("obteniendo hooks_token de la tienda: %w", err)
	}
	if storeInfo.HooksToken == "" {
		return "", domain.ErrWebhookMissingToken
	}

	uc.persistStoreInfo(ctx, integrationID, integration, storeInfo)

	return storeInfo.HooksToken, nil
}

func (uc *jumpsellerUseCase) persistStoreInfo(ctx context.Context, integrationID string, integration *domain.Integration, storeInfo *domain.StoreInfo) {
	config := integration.Config
	if config == nil {
		config = map[string]interface{}{}
	}
	config["hooks_token"] = storeInfo.HooksToken
	config["store_code"] = storeInfo.Code
	config["store_name"] = storeInfo.Name
	config["store_url"] = storeInfo.URL

	if err := uc.service.UpdateIntegrationConfig(ctx, integrationID, config); err != nil {
		uc.logger.Warn(ctx).Err(err).
			Str("integration_id", integrationID).
			Msg("No se pudo persistir la informacion de la tienda Jumpseller")
		return
	}
	integration.Config = config
}

func (uc *jumpsellerUseCase) ProcessWebhookOrder(ctx context.Context, event string, storeCode string, integrationID string, rawBody []byte) error {
	if ignoredWebhookEvents[event] {
		uc.logger.Info(ctx).
			Str("event", event).
			Str("store_code", storeCode).
			Msg("Evento de webhook de Jumpseller ignorado")
		return nil
	}

	if integrationID == "" {
		return fmt.Errorf("webhook sin integration_id en la URL: recrea los webhooks de la integracion para incluirlo")
	}

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("resolviendo integracion %s del webhook: %w", integrationID, err)
	}
	if integration == nil {
		return fmt.Errorf("integracion %s del webhook no existe", integrationID)
	}

	var envelope response.OrderEnvelope
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return fmt.Errorf("deserializing webhook payload: %w", err)
	}

	if envelope.Order.ID == 0 {
		return fmt.Errorf("payload de webhook sin orden valida")
	}

	order := envelope.Order.ToDomain()

	uc.logger.Info(ctx).
		Str("event", event).
		Str("store_code", storeCode).
		Int64("order_id", order.ID).
		Uint("integration_id", integration.ID).
		Interface("business_id", integration.BusinessID).
		Msg("Processing Jumpseller webhook order")

	dto := mapper.MapJumpsellerOrderToProbability(&order, rawBody)
	dto.IntegrationID = integration.ID
	dto.BusinessID = integration.BusinessID

	if err := uc.publisher.Publish(ctx, dto); err != nil {
		uc.logger.Error(ctx).Err(err).
			Str("order_number", dto.OrderNumber).
			Str("event", event).
			Msg("Failed to publish webhook order")
		return fmt.Errorf("publishing webhook order: %w", err)
	}

	uc.logger.Info(ctx).
		Str("order_number", dto.OrderNumber).
		Str("event", event).
		Msg("Jumpseller webhook order published successfully")

	return nil
}
