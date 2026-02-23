package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

// GetIntegrationByExternalID busca una integración por su identificador externo (ej: shop domain) y tipo.
func (uc *IntegrationUseCase) GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*domain.PublicIntegration, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationByExternalID")

	var typeID *uint
	if integrationType > 0 {
		tid := uint(integrationType)
		typeID = &tid
	}

	filters := domain.IntegrationFilters{
		Page:              1,
		PageSize:          1,
		IntegrationTypeID: typeID,
		StoreID:           &externalID,
	}

	integrations, _, err := uc.repo.ListIntegrations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing integrations: %w", err)
	}
	if len(integrations) == 0 {
		return nil, fmt.Errorf("integration not found for external_id %s", externalID)
	}

	return uc.mapToPublicIntegration(integrations[0]), nil
}

// UpdateIntegrationConfig actualiza el config de una integración haciendo merge con el config existente.
func (uc *IntegrationUseCase) UpdateIntegrationConfig(ctx context.Context, integrationID string, newConfig map[string]interface{}) error {
	ctx = log.WithFunctionCtx(ctx, "UpdateIntegrationConfig")

	id, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		return fmt.Errorf("ID de integración inválido: %w", err)
	}

	existing, err := uc.GetPublicIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("error al obtener integración: %w", err)
	}

	// Obtener config existente
	existingConfig := existing.Config
	if existingConfig == nil {
		existingConfig = make(map[string]interface{})
	}

	// Hacer merge
	for k, v := range newConfig {
		existingConfig[k] = v
	}

	// Convertir a JSON
	configBytes, err := json.Marshal(existingConfig)
	if err != nil {
		return fmt.Errorf("error al serializar config: %w", err)
	}

	dto := domain.UpdateIntegrationDTO{
		Config: func() *datatypes.JSON {
			configJSON := datatypes.JSON(configBytes)
			return &configJSON
		}(),
	}

	_, err = uc.UpdateIntegration(ctx, uint(id), dto)
	if err != nil {
		return fmt.Errorf("error al actualizar integración: %w", err)
	}

	return nil
}

// TestConnectionFromConfig prueba la conexión con datos de config/credentials proporcionados directamente.
func (uc *IntegrationUseCase) TestConnectionFromConfig(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	ctx = log.WithFunctionCtx(ctx, "TestConnectionFromConfig")

	// Intentar obtener como int primero, luego como string para compatibilidad
	var integrationType int
	if intVal, ok := config["integration_type"].(int); ok {
		integrationType = intVal
	} else if floatVal, ok := config["integration_type"].(float64); ok {
		integrationType = int(floatVal)
	} else if strVal, ok := config["integration_type"].(string); ok {
		integrationType = domain.IntegrationTypeCodeAsInt(strVal)
	} else {
		return fmt.Errorf("integration_type is required in config and must be int or string")
	}

	if integrationType == 0 {
		return fmt.Errorf("integration_type cannot be zero")
	}

	provider, ok := uc.providerReg.Get(integrationType)
	if !ok {
		return fmt.Errorf("integration no registrada para tipo %d", integrationType)
	}
	return provider.TestConnection(ctx, config, credentials)
}

// OnIntegrationCreated registra un observador que se ejecuta cuando se crea una integración del tipo especificado.
func (uc *IntegrationUseCase) OnIntegrationCreated(integrationType int, observer func(context.Context, *domain.PublicIntegration)) {
	uc.RegisterObserver(func(ctx context.Context, integration *domain.Integration) {
		var integrationTypeCode int
		if integration.IntegrationType != nil {
			integrationTypeCode = domain.IntegrationTypeCodeAsInt(integration.IntegrationType.Code)
		}

		if integrationTypeCode == integrationType {
			publicIntegration := uc.mapToPublicIntegration(integration)
			observer(ctx, publicIntegration)
		}
	})
}
