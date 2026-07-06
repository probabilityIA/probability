package usecases

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/utils"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
)

func (uc *SyncOrdersUseCase) CreateWebhook(ctx context.Context, integrationID string, baseURL string) (*domain.CreateWebhookResult, error) {
	integration, err := uc.integrationService.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener integración: %w", err)
	}

	accessToken, err := uc.integrationService.DecryptCredential(ctx, integrationID, "access_token")
	if err != nil {
		return nil, fmt.Errorf("error al obtener access_token: %w", err)
	}

	storeName, ok := integration.Config["store_name"].(string)
	if !ok || storeName == "" {
		return nil, fmt.Errorf("store_name no encontrado en la configuración")
	}

	storeName = utils.ResolveEffectiveStoreDomain(integration, storeName)

	apiPath := "/api/v1/integrations/shopify/webhook"
	if strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}
	if strings.HasSuffix(baseURL, "/api/v1") {
		baseURL = strings.TrimSuffix(baseURL, "/api/v1")
	}

	webhookURL := fmt.Sprintf("%s%s", baseURL, apiPath)

	uc.log.Info(ctx).
		Str("integration_id", integrationID).
		Str("store_name", storeName).
		Str("webhook_url", webhookURL).
		Bool("is_testing", integration.IsTesting).
		Msg("Creando webhooks en Shopify")

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error al parsear baseURL: %w", err)
	}

	hostname := strings.ToLower(parsedURL.Hostname())
	isPrivateHost := hostname == "localhost" ||
		strings.HasPrefix(hostname, "127.") ||
		hostname == "::1" ||
		hostname == "host.docker.internal" ||
		strings.HasSuffix(hostname, ".local")

	if parsedURL.Scheme != "https" || isPrivateHost {
		return &domain.CreateWebhookResult{
			ExistingWebhooks: []domain.WebhookInfo{},
			DeletedWebhooks:  []domain.WebhookInfo{},
			CreatedWebhooks:  []string{},
			WebhookURL:       webhookURL,
		}, fmt.Errorf("Shopify requiere una URL de webhook HTTPS publica; la actual (%s) no lo es. Configura WEBHOOK_BASE_URL con un dominio HTTPS publico (en local, expon el backend con un tunel como cloudflared)", webhookURL)
	}

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

	for _, webhook := range existingWebhooks {
		if err := uc.shopifyClient.DeleteWebhook(ctx, storeName, accessToken, webhook.ID); err != nil {
			uc.log.Warn(ctx).Err(err).Str("webhook_id", webhook.ID).Str("topic", webhook.Topic).Msg("Error al eliminar webhook existente")
			continue
		}
		result.DeletedWebhooks = append(result.DeletedWebhooks, webhook)
	}

	events := []string{
		"orders/create",
		"orders/updated",
		"orders/paid",
		"orders/cancelled",
		"orders/fulfilled",
	}

	webhookConfigured := true
	var lastErr error
	for _, event := range events {
		webhookID, err := uc.shopifyClient.CreateWebhook(ctx, storeName, accessToken, webhookURL, event)
		if err != nil {
			webhookConfigured = false
			lastErr = err
			uc.log.Error(ctx).Err(err).Str("event", event).Str("webhook_url", webhookURL).Str("store_name", storeName).Msg("Error al crear webhook en Shopify")
			continue
		}
		result.CreatedWebhooks = append(result.CreatedWebhooks, webhookID)
	}

	if len(result.CreatedWebhooks) == 0 {
		return result, webhookCreationError(lastErr)
	}

	configUpdate := map[string]interface{}{
		"webhook_url":        webhookURL,
		"webhook_configured": webhookConfigured,
		"webhook_ids":        result.CreatedWebhooks,
	}

	if err := uc.integrationService.UpdateIntegrationConfig(ctx, integrationID, configUpdate); err != nil {
		return result, fmt.Errorf("error al actualizar config de la integración: %w", err)
	}

	return result, nil
}

func webhookCreationError(shopifyErr error) error {
	detail := ""
	if shopifyErr != nil {
		detail = shopifyErr.Error()
	}
	low := strings.ToLower(detail)

	switch {
	case strings.Contains(low, "protected customer data") ||
		strings.Contains(low, "permission to create or update webhooks"):
		return fmt.Errorf("Shopify bloqueo la creacion de webhooks de ordenes porque la app no tiene aprobado el acceso a los datos protegidos del cliente ('Protected customer data'). " +
			"Los eventos orders/* contienen datos del cliente y Shopify los rechaza (403) mientras la app no tenga ese acceso aprobado, aunque el token sea valido. " +
			"Como habilitarlo: 1) Entra al Shopify Partner Dashboard, abre tu app y ve a 'API access' -> 'Protected customer data access'. " +
			"2) Solicita/aprueba el acceso a 'Customer data' y a 'Orders', completando el cuestionario de proteccion de datos. " +
			"3) Confirma que la app tenga los scopes 'read_orders' y 'write_orders'. " +
			"4) Reinstala o re-autoriza la app para regenerar el Access Token (idealmente token offline 'shpat_', no online 'shpua_' que caduca) y actualiza el token en la integracion. " +
			"Luego reconecta la integracion para volver a crear los webhooks.")
	case strings.Contains(low, "access scope") || strings.Contains(low, "invalid topic"):
		return fmt.Errorf("La app de Shopify no tiene permiso para suscribirse a los eventos de ordenes: le falta el scope 'read_orders', que Shopify trata como datos protegidos del cliente. " +
			"Como habilitarlo: 1) En el admin de Shopify entra a Configuracion -> Apps y canales de venta -> Desarrollar apps y abre tu Custom App. " +
			"2) En Configuration -> Admin API integration, activa los scopes 'read_orders' y 'write_orders' (y 'read_fulfillments' si registras eventos de fulfillment). " +
			"3) Guarda y aprueba el acceso a 'Protected customer data'. " +
			"4) Reinstala o re-autoriza la app para regenerar el Access Token con los nuevos permisos y actualiza el token en la integracion. " +
			"Luego vuelve a crear los webhooks.")
	case strings.Contains(low, "token de acceso") || strings.Contains(low, "unauthorized") || strings.Contains(low, "expir"):
		return fmt.Errorf("El token de acceso de Shopify es invalido o expiro. Reconecta la integracion para regenerar el token y vuelve a crear los webhooks.")
	case strings.Contains(low, "https") || strings.Contains(low, "protocol http"):
		return fmt.Errorf("Shopify rechazo la URL del webhook porque no es HTTPS publica. Configura WEBHOOK_BASE_URL con un dominio HTTPS publico (en produccion ya lo es; en local usa un tunel como cloudflared).")
	case detail != "":
		return fmt.Errorf("no se pudo crear ningun webhook en Shopify. Detalle: %s", detail)
	default:
		return fmt.Errorf("no se pudo crear ningun webhook en Shopify")
	}
}
