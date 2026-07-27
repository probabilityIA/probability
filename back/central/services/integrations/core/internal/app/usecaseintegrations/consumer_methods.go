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

func (uc *IntegrationUseCase) GetIntegrationByExternalID(ctx context.Context, externalID string, integrationType int) (*domain.PublicIntegration, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationByExternalID")

	if integrationType > 0 && externalID != "" {
		if cached, err := uc.cache.GetByStoreAndType(ctx, externalID, uint(integrationType)); err == nil && cached != nil && cached.IsActive {
			uc.log.Debug(ctx).Str("external_id", externalID).Msg("Cache hit - store index")
			return uc.mapToPublicIntegration(cachedToIntegration(cached)), nil
		}
	}

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

	uc.cacheIntegrationMeta(ctx, integrations[0])

	return uc.mapToPublicIntegration(integrations[0]), nil
}

func cachedToIntegration(cached *domain.CachedIntegration) *domain.Integration {
	configJSON, _ := json.Marshal(cached.Config)
	return &domain.Integration{
		ID:                cached.ID,
		Name:              cached.Name,
		Code:              cached.Code,
		Category:          cached.Category,
		IntegrationTypeID: cached.IntegrationTypeID,
		BusinessID:        cached.BusinessID,
		StoreID:           cached.StoreID,
		IsActive:          cached.IsActive,
		IsDefault:         cached.IsDefault,
		IsTesting:         cached.IsTesting,
		Config:            datatypes.JSON(configJSON),
		Description:       cached.Description,
		CreatedAt:         cached.CreatedAt,
		UpdatedAt:         cached.UpdatedAt,
		IntegrationType: &domain.IntegrationType{
			ID:          cached.IntegrationTypeID,
			Code:        cached.IntegrationTypeCode,
			BaseURL:     cached.BaseURL,
			BaseURLTest: cached.BaseURLTest,
		},
	}
}

func (uc *IntegrationUseCase) cacheIntegrationMeta(ctx context.Context, integration *domain.Integration) {
	configMap := make(map[string]interface{})
	if len(integration.Config) > 0 {
		json.Unmarshal(integration.Config, &configMap)
	}

	integrationType := integration.IntegrationType
	if integrationType == nil && integration.IntegrationTypeID > 0 {
		loaded, err := uc.repo.GetIntegrationTypeByID(ctx, integration.IntegrationTypeID)
		if err != nil {
			uc.log.Warn(ctx).Err(err).
				Uint("integration_id", integration.ID).
				Uint("integration_type_id", integration.IntegrationTypeID).
				Msg("No se pudo cargar el tipo de integracion para cachear base_url/base_url_test")
		}
		integrationType = loaded
	}

	integrationTypeCode := ""
	baseURL := ""
	baseURLTest := ""
	if integrationType != nil {
		integrationTypeCode = integrationType.Code
		baseURL = integrationType.BaseURL
		baseURLTest = integrationType.BaseURLTest
	}

	cachedMeta := &domain.CachedIntegration{
		ID:                  integration.ID,
		Name:                integration.Name,
		Code:                integration.Code,
		Category:            integration.Category,
		IntegrationTypeID:   integration.IntegrationTypeID,
		IntegrationTypeCode: integrationTypeCode,
		BusinessID:          integration.BusinessID,
		StoreID:             integration.StoreID,
		IsActive:            integration.IsActive,
		IsDefault:           integration.IsDefault,
		IsTesting:           integration.IsTesting,
		Config:              configMap,
		Description:         integration.Description,
		CreatedAt:           integration.CreatedAt,
		UpdatedAt:           integration.UpdatedAt,
		BaseURL:             baseURL,
		BaseURLTest:         baseURLTest,
	}

	if err := uc.cache.SetIntegration(ctx, cachedMeta); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to cache integration metadata")
	}
}

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

	existingConfig := existing.Config
	if existingConfig == nil {
		existingConfig = make(map[string]interface{})
	}

	for k, v := range newConfig {
		existingConfig[k] = v
	}

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

func (uc *IntegrationUseCase) TestConnectionFromConfig(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	ctx = log.WithFunctionCtx(ctx, "TestConnectionFromConfig")

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
