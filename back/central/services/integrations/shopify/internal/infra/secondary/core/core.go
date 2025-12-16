package core

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/domain"
)

// integrationServiceAdapter adapta core.IIntegrationCore a domain.IIntegrationService
// Esto permite desacoplar el dominio de módulos externos (arquitectura hexagonal)
type integrationServiceAdapter struct {
	coreIntegration core.IIntegrationCore
}

func (a *integrationServiceAdapter) GetIntegrationByID(ctx context.Context, integrationID string) (*domain.Integration, error) {
	integration, err := a.coreIntegration.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	return &domain.Integration{
		ID:              integration.ID,
		BusinessID:      integration.BusinessID,
		Name:            integration.Name,
		StoreID:         integration.StoreID,
		IntegrationType: integration.IntegrationType,
		Config:          integration.Config,
	}, nil
}

func (a *integrationServiceAdapter) GetIntegrationByStoreID(ctx context.Context, storeID string) (*domain.Integration, error) {
	integration, err := a.coreIntegration.GetIntegrationByStoreID(ctx, storeID, core.IntegrationTypeShopify)
	if err != nil {
		return nil, err
	}

	return &domain.Integration{
		ID:              integration.ID,
		BusinessID:      integration.BusinessID,
		Name:            integration.Name,
		StoreID:         integration.StoreID,
		IntegrationType: integration.IntegrationType,
		Config:          integration.Config,
	}, nil
}

func (a *integrationServiceAdapter) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return a.coreIntegration.DecryptCredential(ctx, integrationID, fieldName)
}

func (a *integrationServiceAdapter) UpdateIntegrationConfig(ctx context.Context, integrationID string, config map[string]interface{}) error {
	return a.coreIntegration.UpdateIntegrationConfig(ctx, integrationID, config)
}

type ShopifyCore struct {
	coreIntegration core.IIntegrationCore
	useCase         usecases.IShopifyUseCase
	client          domain.ShopifyClient
}

func New(coreIntegration core.IIntegrationCore, shopifyClient domain.ShopifyClient, orderPublisher domain.OrderPublisher) *ShopifyCore {
	// Crear adaptador que implementa domain.IIntegrationService
	integrationService := &integrationServiceAdapter{
		coreIntegration: coreIntegration,
	}

	useCase := usecases.New(integrationService, shopifyClient, orderPublisher)
	return &ShopifyCore{
		coreIntegration: coreIntegration,
		useCase:         useCase,
		client:          shopifyClient,
	}
}

func (s *ShopifyCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	storeName, ok := config["store_name"].(string)
	if !ok || storeName == "" {
		return fmt.Errorf("el nombre de la tienda (store_name) es requerido")
	}

	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("el token de acceso (access_token) es requerido")
	}

	valid, _, err := s.client.ValidateToken(ctx, storeName, accessToken)
	if err != nil {
		return err // El error ya viene con mensaje descriptivo en español
	}

	if !valid {
		return fmt.Errorf("credenciales o nombre de tienda inválidos")
	}

	return nil
}

func (s *ShopifyCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return s.useCase.SyncOrders(ctx, integrationID)
}

func (s *ShopifyCore) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	return fmt.Errorf("SyncOrdersByBusiness should be handled by core, not by individual syncers")
}

// GetWebhookURL construye la URL del webhook para Shopify
func (s *ShopifyCore) GetWebhookURL(ctx context.Context, baseURL string, integrationID uint) (*core.WebhookInfo, error) {
	// Construir la URL del webhook
	// El webhook se recibe en: /api/v1/integrations/shopify/webhook
	webhookURL := fmt.Sprintf("%s/integrations/shopify/webhook", baseURL)

	return &core.WebhookInfo{
		URL:         webhookURL,
		Method:      "POST",
		Description: "URL para configurar en Shopify Admin > Settings > Notifications > Webhooks. Configure este webhook para recibir eventos de órdenes en tiempo real.",
		Events: []string{
			"orders/create",
			"orders/updated",
			"orders/paid",
			"orders/cancelled",
			"orders/fulfilled",
			"orders/partially_fulfilled",
		},
	}, nil
}

// ListWebhooks lista todos los webhooks de una integración de Shopify
func (s *ShopifyCore) ListWebhooks(ctx context.Context, integrationID string) ([]interface{}, error) {
	webhooks, err := s.useCase.ListWebhooks(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	// Convertir a []interface{} para la interfaz
	result := make([]interface{}, len(webhooks))
	for i, wh := range webhooks {
		result[i] = wh
	}

	return result, nil
}

// DeleteWebhook elimina un webhook de Shopify
func (s *ShopifyCore) DeleteWebhook(ctx context.Context, integrationID, webhookID string) error {
	return s.useCase.DeleteWebhook(ctx, integrationID, webhookID)
}
