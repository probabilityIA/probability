package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/datatypes"
)

// CreateIntegration crea una nueva integración
func (uc *IntegrationUseCase) CreateIntegration(ctx context.Context, dto domain.CreateIntegrationDTO) (*domain.Integration, error) {
	ctx = log.WithFunctionCtx(ctx, "CreateIntegration")

	// Validaciones
	if dto.Name == "" {
		return nil, domain.ErrIntegrationNameRequired
	}
	if dto.Code == "" {
		return nil, domain.ErrIntegrationCodeRequired
	}
	if dto.IntegrationTypeID == 0 {
		return nil, domain.ErrIntegrationTypeRequired
	}
	// Category ya no se valida aquí, se obtiene del IntegrationType

	// Validar que el tipo de integración exista (necesitamos el repositorio de tipos)
	// TODO: Inyectar IIntegrationTypeRepository en el use case

	// Validar que no exista otra integración con el mismo código
	exists, err := uc.repo.ExistsIntegrationByCode(ctx, dto.Code, dto.BusinessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al verificar existencia de código")
		return nil, fmt.Errorf("error al verificar código: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("%w: %s", domain.ErrIntegrationCodeExists, dto.Code)
	}

	// TODO: Validar reglas específicas del tipo de integración
	// Por ejemplo, si el tipo es WhatsApp, debe ser global (BusinessID = NULL)
	// Esto se puede hacer consultando el IntegrationType y sus reglas

	// Obtener el tipo de integración para validar la conexión y derivar la categoría
	integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, dto.IntegrationTypeID)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("integration_type_id", dto.IntegrationTypeID).Msg("Error al obtener tipo de integración")
		return nil, fmt.Errorf("error al obtener tipo de integración: %w", err)
	}

	// Derivar la categoría del IntegrationType
	// Si el IntegrationType tiene una categoría cargada, usar su código
	categoryCode := ""
	if integrationType.Category != nil {
		categoryCode = integrationType.Category.Code
	} else if integrationType.CategoryID > 0 {
		// Si no está cargada pero tenemos el ID, obtenerla
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

	// VALIDAR CONEXIÓN ANTES DE GUARDAR
	// Obtener provider registrado para este tipo
	integrationTypeInt := domain.IntegrationTypeCodeAsInt(integrationType.Code)
	provider, hasProvider := uc.providerReg.Get(integrationTypeInt)
	if !hasProvider {
		uc.log.Warn(ctx).
			Str("type_code", integrationType.Code).
			Msg("No hay provider registrado, solo validando credenciales básicas")

		// Verificar si usa token de plataforma (en ese caso no hay credenciales propias que validar)
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
			// Fallback: validación básica si no hay provider y no usa token de plataforma
			if err := uc.validateBasicCredentials(ctx, integrationType.Code, dto.Credentials); err != nil {
				return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationTestFailed, err)
			}
		}
	} else {
		// Deserializar Config a map para el provider
		var configMap map[string]interface{}
		if len(dto.Config) > 0 {
			if err := json.Unmarshal(dto.Config, &configMap); err != nil {
				return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationConfigDeserialize, err)
			}
		}
		if configMap == nil {
			configMap = make(map[string]interface{})
		}

		// Inyectar base_url del tipo de integración si no viene en el config del usuario
		// El usuario no ingresa estas URLs — vienen de integration_types table
		if _, has := configMap["base_url"]; !has && integrationType.BaseURL != "" {
			configMap["base_url"] = integrationType.BaseURL
		}
		if _, has := configMap["base_url_test"]; !has && integrationType.BaseURLTest != "" {
			configMap["base_url_test"] = integrationType.BaseURLTest
		}

		// Si usa token de plataforma, omitir test de conexión (usa credenciales del tipo)
		usePlatformToken, _ := configMap["use_platform_token"].(bool)
		if usePlatformToken {
			uc.log.Info(ctx).
				Str("type_code", integrationType.Code).
				Str("integration_code", dto.Code).
				Msg("use_platform_token=true, omitiendo test de conexión con provider")
		} else {
			// Testear conexión con el provider específico
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

	// Convertir Config a datatypes.JSON
	var configJSON datatypes.JSON
	if len(dto.Config) > 0 {
		configBytes, err := json.Marshal(dto.Config)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationConfigSerialize, err)
		}
		configJSON = configBytes
	}

	// Convertir Credentials a []byte (se encriptará en el repository)
	var credentialsJSON datatypes.JSON
	if len(dto.Credentials) > 0 {
		// DEBUG: Log credential keys being saved
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

	// Crear entidad de dominio
	integration := &domain.Integration{
		Name:              dto.Name,
		Code:              dto.Code,
		Category:          categoryCode, // Derivado de IntegrationType.Category.Code
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

	// Guardar en repository (encriptará credenciales automáticamente)
	if err := uc.repo.CreateIntegration(ctx, integration); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Error al crear integración")
		return nil, fmt.Errorf("error al crear integración: %w", err)
	}

	// Los observers de OnIntegrationCreated filtran por integration.IntegrationType.Code;
	// la entidad recien creada no trae la relacion cargada, sin esto nunca se disparan.
	integration.IntegrationType = integrationType

	// ✅ NUEVO - Cachear metadata
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

	// ✅ NUEVO - Cachear credentials desencriptadas
	if len(dto.Credentials) > 0 {
		cachedCreds := &domain.CachedCredentials{
			IntegrationID: integration.ID,
			Credentials:   dto.Credentials, // Ya están desencriptadas en el DTO
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

	// Crear webhooks automáticamente si el provider lo soporta
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

	// Notificar observadores (e.g., para auto-sync)
	// Hacemos esto de forma asíncrona para no bloquear la respuesta HTTP
	go func() {
		for _, observer := range uc.observers {
			// Crear un nuevo contexto desconectado del request HTTP cancelar
			bgCtx := context.Background()
			// Tratar pánicos en observadores para no romper nada
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
