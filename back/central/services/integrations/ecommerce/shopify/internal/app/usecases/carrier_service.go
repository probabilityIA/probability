package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/utils"
)

const carrierServiceName = "Probability"

func buildShippingRatesCallbackURL(baseURL, integrationID string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	baseURL = strings.TrimSuffix(baseURL, "/api/v1")
	return fmt.Sprintf("%s/api/v1/shopify/shipping-rates/%s", baseURL, integrationID)
}

func (uc *SyncOrdersUseCase) EnableCarrierCalculatedShipping(ctx context.Context, integrationID string, publicBaseURL string) (string, error) {
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return "", fmt.Errorf("error al obtener integracion: %w", err)
	}

	accessToken, err := uc.integrationService.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return "", fmt.Errorf("error al obtener access_token: %w", err)
	}

	storeName, ok := integration.Config["store_name"].(string)
	if !ok || storeName == "" {
		return "", fmt.Errorf("store_name no encontrado en la configuracion")
	}
	storeName = utils.ResolveEffectiveStoreDomain(integration, storeName)

	if publicBaseURL == "" {
		return "", fmt.Errorf("no hay una URL base publica configurada (WEBHOOK_BASE_URL)")
	}

	callbackURL := buildShippingRatesCallbackURL(publicBaseURL, integrationID)

	if existingID, _ := integration.Config["carrier_service_id"].(string); existingID != "" {
		_ = uc.shopifyClient.DeleteCarrierService(ctx, storeName, accessToken, existingID)
	}

	uc.log.Info(ctx).
		Str("integration_id", integrationID).
		Str("store_name", storeName).
		Str("callback_url", callbackURL).
		Msg("Registrando carrier service en Shopify")

	carrierServiceID, err := uc.shopifyClient.CreateCarrierService(ctx, storeName, accessToken, callbackURL, carrierServiceName)
	if err != nil {
		return "", err
	}

	configUpdate := map[string]interface{}{
		"carrier_calculated_shipping_enabled": true,
		"carrier_service_id":                  carrierServiceID,
		"carrier_service_callback_url":        callbackURL,
	}
	if err := uc.integrationService.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
		return carrierServiceID, fmt.Errorf("carrier service creado pero fallo al guardar el estado: %w", err)
	}

	return carrierServiceID, nil
}

func (uc *SyncOrdersUseCase) DisableCarrierCalculatedShipping(ctx context.Context, integrationID string) error {
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("error al obtener integracion: %w", err)
	}

	accessToken, err := uc.integrationService.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return fmt.Errorf("error al obtener access_token: %w", err)
	}

	storeName, ok := integration.Config["store_name"].(string)
	if !ok || storeName == "" {
		return fmt.Errorf("store_name no encontrado en la configuracion")
	}
	storeName = utils.ResolveEffectiveStoreDomain(integration, storeName)

	carrierServiceID, _ := integration.Config["carrier_service_id"].(string)
	if carrierServiceID != "" {
		if err := uc.shopifyClient.DeleteCarrierService(ctx, storeName, accessToken, carrierServiceID); err != nil {
			return err
		}
	}

	configUpdate := map[string]interface{}{
		"carrier_calculated_shipping_enabled": false,
		"carrier_service_id":                  "",
		"carrier_service_callback_url":        "",
	}
	if err := uc.integrationService.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
		return fmt.Errorf("carrier service eliminado pero fallo al guardar el estado: %w", err)
	}

	return nil
}
