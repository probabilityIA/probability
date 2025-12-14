package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

const (
	IntegrationTypeShopify      = 1
	IntegrationTypeWhatsApp     = 2
	IntegrationTypeMercadoLibre = 3
)

type IntegrationWithCredentials = domain.IntegrationWithCredentials

type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          interface{}
}

func (ic *integrationCore) GetIntegrationByID(ctx context.Context, integrationID string) (*Integration, error) {
	publicIntegration, err := ic.useCase.GetPublicIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	return &Integration{
		ID:              publicIntegration.ID,
		BusinessID:      publicIntegration.BusinessID,
		Name:            publicIntegration.Name,
		StoreID:         publicIntegration.StoreID,
		IntegrationType: publicIntegration.IntegrationType,
		Config:          publicIntegration.Config,
	}, nil
}

// GetIntegrationByStoreID busca una integración por StoreID (ej: shop domain) y tipo
func (ic *integrationCore) GetIntegrationByStoreID(ctx context.Context, storeID string, integrationType int) (*Integration, error) {
	var typeID *uint
	if integrationType > 0 {
		tid := uint(integrationType)
		typeID = &tid
	}

	filters := domain.IntegrationFilters{
		Page:              1,
		PageSize:          1,
		IntegrationTypeID: typeID,
		StoreID:           &storeID,
	}

	integrations, _, err := ic.useCase.ListIntegrations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing integrations: %w", err)
	}
	if len(integrations) == 0 {
		return nil, fmt.Errorf("integration not found for store_id %s", storeID)
	}

	integration := integrations[0]

	var config map[string]interface{}
	if len(integration.Config) > 0 {
		_ = json.Unmarshal(integration.Config, &config)
	}

	integrationTypeCode := integrationType
	if integrationTypeCode == 0 && integration.IntegrationType != nil {
		integrationTypeCode = getIntegrationTypeCodeAsInt(integration.IntegrationType.Code)
	} else if integrationTypeCode == 0 {
		integrationTypeCode = int(integration.IntegrationTypeID)
	}

	return &Integration{
		ID:              integration.ID,
		BusinessID:      integration.BusinessID,
		Name:            integration.Name,
		StoreID:         integration.StoreID,
		IntegrationType: integrationTypeCode,
		Config:          config,
	}, nil
}

func (ic *integrationCore) DecryptCredential(ctx context.Context, integrationID string, fieldName string) (string, error) {
	return ic.useCase.DecryptCredentialField(ctx, integrationID, fieldName)
}

func (ic *integrationCore) RegisterIntegration(integrationType int, integration IIntegrationContract) {
	if integrationType == 0 {
		ic.logger.Error().Msg("RegisterIntegration: integration type cannot be zero")
		return
	}
	if integration == nil {
		ic.logger.Error().Msg("RegisterIntegration: integration cannot be nil")
		return
	}

	ic.integrations[integrationType] = integration

	useCaseImpl, ok := ic.useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		ic.logger.Error().Msg("RegisterIntegration: error interno: no se puede acceder al registry de testers")
		return
	}

	adapter := &integrationAdapter{integration: integration}
	if err := useCaseImpl.GetTesterRegistry().Register(integrationType, adapter); err != nil {
		ic.logger.Error().Err(err).Int("integration_type", integrationType).Msg("RegisterIntegration: error al registrar tester")
		return
	}

	ic.logger.Info().Int("integration_type", integrationType).Msg("Integration registered successfully")
}

func (ic *integrationCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	// Intentar obtener como int primero, luego como string para compatibilidad
	var integrationType int
	if intVal, ok := config["integration_type"].(int); ok {
		integrationType = intVal
	} else if floatVal, ok := config["integration_type"].(float64); ok {
		// JSON numbers se deserializan como float64
		integrationType = int(floatVal)
	} else if strVal, ok := config["integration_type"].(string); ok {
		// Compatibilidad con strings antiguos
		integrationType = getIntegrationTypeCodeAsInt(strVal)
	} else {
		return fmt.Errorf("integration_type is required in config and must be int or string")
	}

	if integrationType == 0 {
		return fmt.Errorf("integration_type cannot be zero")
	}

	integration, ok := ic.integrations[integrationType]
	if !ok {
		return fmt.Errorf("integration no registrada para tipo %d", integrationType)
	}
	return integration.TestConnection(ctx, config, credentials)
}

func (ic *integrationCore) SyncOrdersByIntegrationID(ctx context.Context, integrationID string) error {
	integration, err := ic.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return err
	}

	integrationImpl, ok := ic.integrations[integration.IntegrationType]
	if !ok {
		return fmt.Errorf("integration no registrada para tipo %d", integration.IntegrationType)
	}

	return integrationImpl.SyncOrdersByIntegrationID(ctx, integrationID)
}

func (ic *integrationCore) SyncOrdersByBusiness(ctx context.Context, businessID uint) error {
	useCaseImpl, ok := ic.useCase.(*usecaseintegrations.IntegrationUseCase)
	if !ok {
		return fmt.Errorf("error interno: no se puede acceder al use case")
	}

	businessIDPtr := &businessID
	filters := domain.IntegrationFilters{
		BusinessID: businessIDPtr,
		IsActive:   &[]bool{true}[0],
	}

	integrations, _, err := useCaseImpl.ListIntegrations(ctx, filters)
	if err != nil {
		return fmt.Errorf("error al obtener integraciones: %w", err)
	}

	for _, integration := range integrations {
		if integration.IntegrationType == nil {
			continue
		}

		integrationID := fmt.Sprintf("%d", integration.ID)
		if err := ic.SyncOrdersByIntegrationID(ctx, integrationID); err != nil {
			continue
		}
	}

	return nil
}

func (ic *integrationCore) RegisterObserverForType(integrationType int, observer func(context.Context, *Integration)) {
	ic.useCase.RegisterObserver(func(ctx context.Context, integration *domain.Integration) {
		var integrationTypeCode int
		if integration.IntegrationType != nil {
			// Convertir el código del tipo de integración a int si es necesario
			// Por ahora asumimos que el código puede ser convertido o comparado
			integrationTypeCode = getIntegrationTypeCodeAsInt(integration.IntegrationType.Code)
		}

		if integrationTypeCode == integrationType {
			publicIntegration := mapDomainToPublicIntegration(ic.useCase, integration)
			observer(ctx, publicIntegration)
		}
	})
}

// getIntegrationTypeCodeAsInt convierte el código de tipo de integración a int
// Esta función mapea los códigos antiguos (strings) a los nuevos (int)
func getIntegrationTypeCodeAsInt(code string) int {
	switch code {
	case "shopify":
		return IntegrationTypeShopify
	case "whatsapp":
		return IntegrationTypeWhatsApp
	case "mercado_libre":
		return IntegrationTypeMercadoLibre
	default:
		return 0
	}
}

type integrationAdapter struct {
	integration IIntegrationContract
}

func (a *integrationAdapter) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return a.integration.TestConnection(ctx, config, credentials)
}
