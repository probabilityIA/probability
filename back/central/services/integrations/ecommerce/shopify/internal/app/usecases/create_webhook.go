package usecases

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases/utils"
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
	storeName, ok := integration.Config["store_name"].(string)
	if !ok || storeName == "" {
		return nil, fmt.Errorf("store_name no encontrado en la configuración")
	}

	// En modo test, usar la URL de pruebas
	storeName = utils.ResolveEffectiveStoreDomain(integration, storeName)

	// Construir nuestra URL del webhook
	// Asegurar que usamos el prefijo /api/v1 ya que el router lo espera
	// Si baseURL ya tiene /api/v1, lo manejamos (aunque asumimos que es el host base)
	apiPath := "/api/v1/integrations/shopify/webhook"
	if strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}
	// Si el usuario configuró baseURL con /api/v1, evitamos duplicarlo
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

	// Validar si la URL es localhost (entorno de pruebas)
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
			uc.log.Warn(ctx).Err(err).Str("webhook_id", webhook.ID).Str("topic", webhook.Topic).Msg("Error al eliminar webhook existente")
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
	}

	// Crear webhooks para todos los eventos
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

	// Si no se creó ningún webhook, retornar un error claro y accionable
	if len(result.CreatedWebhooks) == 0 {
		return result, webhookCreationError(lastErr)
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

// webhookCreationError construye un mensaje claro y accionable cuando ningun
// webhook pudo crearse, detectando las causas mas comunes (scope faltante en la
// app de Shopify, token invalido) y explicando como resolverlas.
func webhookCreationError(shopifyErr error) error {
	detail := ""
	if shopifyErr != nil {
		detail = shopifyErr.Error()
	}
	low := strings.ToLower(detail)

	switch {
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
