package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

func (uc *IntegrationUseCase) CreateIntegration(ctx context.Context, dto domain.CreateIntegrationDTO) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "CreateIntegration")

	if dto.Name == "" {
		return nil, domain.ErrIntegrationNameRequired
	}
	if dto.Code == "" {
		return nil, domain.ErrIntegrationCodeRequired
	}
	if dto.IntegrationTypeID == 0 {
		return nil, domain.ErrIntegrationTypeRequired
	}

	exists, err := uc.repo.ExistsIntegrationByCode(ctx, dto.Code, dto.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al verificar existencia de código")
		return nil, fmt.Errorf("error al verificar código: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationCodeExists, dto.Code)
	}

	if err := uc.ensureStoreIDNotInUse(ctx, dto.StoreID, dto.IntegrationTypeID, dto.BusinessID, 0); err != nil {
		return nil, err
	}

	integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, dto.IntegrationTypeID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_type_id", dto.IntegrationTypeID).Msg("Error al obtener tipo de integración")
		return nil, fmt.Errorf("error al obtener tipo de integración: %w", err)
	}

	categoryCode := ""
	if integrationType.Category != nil {
		categoryCode = integrationType.Category.Code
	} else if integrationType.CategoryID > 0 {
		category, err := uc.repo.GetIntegrationCategoryByID(ctx, integrationType.CategoryID)
		if err != nil {
			uc.log.Error(ctx).Err(err).Uint("category_id", integrationType.CategoryID).Msg("Error al obtener categoría de integración")
			return nil, fmt.Errorf("error al obtener categoría: %w", err)
		}
		categoryCode = category.Code
	}

	if categoryCode == "" {
		uc.log.Error(ctx).Uint("integration_type_id", dto.IntegrationTypeID).Msg("Tipo de integración sin categoría válida")
		return nil, domain.ErrIntegrationCategoryInvalid
	}

	if categoryCode == "ecommerce" && dto.BusinessID != nil && uc.ecommerceLimitChecker != nil {
		limit, err := uc.ecommerceLimitChecker(ctx, *dto.BusinessID)
		if err == nil && limit > 0 {
			active := true
			category := categoryCode
			_, total, err := uc.repo.ListIntegrations(ctx, domain.IntegrationFilters{
				BusinessID: dto.BusinessID,
				Category:   &category,
				IsActive:   &active,
			})
			if err == nil && total >= int64(limit) {
				return nil, domain.ErrEcommerceLimitReached
			}
		}
	}

	integrationTypeInt := domain.IntegrationTypeCodeAsInt(integrationType.Code)
	provider, hasProvider := uc.providerReg.Get(integrationTypeInt)
	if !hasProvider {
		uc.log.Warn(ctx).
			Str("type_code", integrationType.Code).
			Msg("No hay provider registrado, solo validando credenciales básicas")

		var configMap map[string]interface{}
		if len(dto.Config) > 0 {
			json.Unmarshal(dto.Config, &configMap) //nolint:errcheck
		}
		usePlatformToken, _ := configMap["use_platform_token"].(bool)

		if usePlatformToken {
			uc.log.Info(ctx).
				Str("type_code", integrationType.Code).
				Msg("use_platform_token=true, omitiendo validación de credenciales propias")
		} else {
			if err := uc.validateBasicCredentials(ctx, integrationType.Code, dto.Credentials); err != nil {
				return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTestFailed, err)
			}
		}
	} else {
		var configMap map[string]interface{}
		if len(dto.Config) > 0 {
			if err := json.Unmarshal(dto.Config, &configMap); err != nil {
				return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationConfigDeserialize, err)
			}
		}
		if configMap == nil {
			configMap = make(map[string]interface{})
		}

		if _, has := configMap["base_url"]; !has && integrationType.BaseURL != "" {
			configMap["base_url"] = integrationType.BaseURL
		}
		if _, has := configMap["base_url_test"]; !has && integrationType.BaseURLTest != "" {
			configMap["base_url_test"] = integrationType.BaseURLTest
		}

		usePlatformToken, _ := configMap["use_platform_token"].(bool)
		if usePlatformToken {
			uc.log.Info(ctx).
				Str("type_code", integrationType.Code).
				Str("integration_code", dto.Code).
				Msg("use_platform_token=true, omitiendo test de conexión con provider")
		} else {
			if err := provider.TestConnection(ctx, configMap, dto.Credentials); err != nil {
				uc.log.Error(ctx).
					Err(err).
					Str("type_code", integrationType.Code).
					Str("integration_code", dto.Code).
					Msg("Test de conexión falló al crear integración")
				return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTestFailed, err)
			}
			uc.log.Info(ctx).
				Str("type_code", integrationType.Code).
				Str("integration_code", dto.Code).
				Msg("Test de conexión exitoso antes de crear integración")
		}
	}

	var configJSON datatypes.JSON
	if len(dto.Config) > 0 {
		configBytes, err := json.Marshal(dto.Config)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationConfigSerialize, err)
		}
		configJSON = configBytes
	}

	var credentialsJSON datatypes.JSON
	if len(dto.Credentials) > 0 {
		credKeys := make([]string, 0, len(dto.Credentials))
		for k := range dto.Credentials {
			credKeys = append(credKeys, k)
		}
		uc.log.Debug(ctx).
			Strs("credential_keys", credKeys).
			Msg("Credentials keys being saved to integration")

		credentialsBytes, err := json.Marshal(dto.Credentials)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsSerialize, err)
		}
		credentialsJSON = credentialsBytes
	}

	integration := &domain.Integration{
		Name:              dto.Name,
		Code:              dto.Code,
		Category:          categoryCode,
		IntegrationTypeID: dto.IntegrationTypeID,
		BusinessID:        dto.BusinessID,
		StoreID:           dto.StoreID,
		IsActive:          dto.IsActive,
		IsDefault:         dto.IsDefault,
		IsTesting:         dto.IsTesting,
		Config:            configJSON,
		Credentials:       credentialsJSON,
		Description:       dto.Description,
		CreatedByID:       dto.CreatedByID,
	}

	if err := uc.repo.CreateIntegration(ctx, integration); err != nil {
		if isUniqueStoreIDViolation(err) {
			uc.log.Warn(ctx).Str("store_id", dto.StoreID).Msg("La base rechazo una cuenta de canal duplicada")
			return nil, domain.ErrIntegrationStoreIDInUse
		}
		uc.log.Error(ctx).Err(err).Msg("Error al crear integración")
		return nil, fmt.Errorf("error al crear integración: %w", err)
	}

	integration.IntegrationType = integrationType

	configMap := make(map[string]interface{})
	if len(integration.Config) > 0 {
		json.Unmarshal(integration.Config, &configMap)
	}

	cachedMeta := &domain.CachedIntegration{
		ID:                  integration.ID,
		Name:                integration.Name,
		Code:                integration.Code,
		Category:            integration.Category,
		IntegrationTypeID:   integration.IntegrationTypeID,
		IntegrationTypeCode: integrationType.Code,
		BusinessID:          integration.BusinessID,
		StoreID:             integration.StoreID,
		IsActive:            integration.IsActive,
		IsDefault:           integration.IsDefault,
		IsTesting:           integration.IsTesting,
		Config:              configMap,
		Description:         integration.Description,
		CreatedAt:           integration.CreatedAt,
		UpdatedAt:           integration.UpdatedAt,
	}

	if err := uc.cache.SetIntegration(ctx, cachedMeta); err != nil {
		uc.log.Warn(ctx).Err(err).Msg("Failed to cache integration metadata")
	}

	if len(dto.Credentials) > 0 {
		cachedCreds := &domain.CachedCredentials{
			IntegrationID: integration.ID,
			Credentials:   dto.Credentials,
		}

		if err := uc.cache.SetCredentials(ctx, cachedCreds); err != nil {
			uc.log.Warn(ctx).Err(err).Msg("Failed to cache credentials")
		}
	}

	uc.log.Info(ctx).
		Uint("id", integration.ID).
		Uint("integration_type_id", integration.IntegrationTypeID).
		Str("code", integration.Code).
		Str("code", integration.Code).
		Msg("Integración creada exitosamente")

	if hasProvider {
		integrationIDStr := fmt.Sprintf("%d", integration.ID)

		go func() {
			bgCtx := context.Background()

			uc.log.Info(bgCtx).
				Str("integration_id", integrationIDStr).
				Msg("🔄 Iniciando creación automática de webhooks...")

			if _, err := uc.CreateWebhookForIntegration(bgCtx, integrationIDStr); err != nil {
				uc.log.Error(bgCtx).
					Err(err).
					Str("integration_id", integrationIDStr).
					Msg("⚠️ Falló la creación automática de webhooks (se puede reintentar manualmente)")
			} else {
				uc.log.Info(bgCtx).
					Str("integration_id", integrationIDStr).
					Msg("✅ Webhooks creados automáticamente exitosamente")
			}
		}()
	}

	go func() {
		for _, observer := range uc.observers {
			bgCtx := context.Background()
			defer func() {
				if r := recover(); r != nil {
					uc.log.Error(bgCtx).Interface("recover", r).Msg("Pánico en observador de creación de integración")
				}
			}()
			observer(bgCtx, integration)
		}
	}()

	return integration, nil
}
