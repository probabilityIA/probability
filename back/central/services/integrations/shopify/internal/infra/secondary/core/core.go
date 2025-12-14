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
		return fmt.Errorf("store_name is required in config")
	}

	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("access_token is required in credentials")
	}

	valid, _, err := s.client.ValidateToken(ctx, storeName, accessToken)
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}

	if !valid {
		return fmt.Errorf("invalid credentials or store name")
	}

	return nil
}

func (s *ShopifyCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	return s.useCase.SyncOrders(ctx, integrationID)
}

func (s *ShopifyCore) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	return fmt.Errorf("SyncOrdersByBusiness should be handled by core, not by individual syncers")
}
