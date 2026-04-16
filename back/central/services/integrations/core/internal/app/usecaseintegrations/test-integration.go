package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// TestIntegration prueba la conexión de una integración usando el provider registrado
func (uc *IntegrationUseCase) TestIntegration(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "TestIntegration")

	// Obtener integración
	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrIntegrationNotFound, err)
	}

	// Convertir Config a map
	var configMap map[string]interface{}
	if len(integration.Config) > 0 {
		if err := json.Unmarshal(integration.Config, &configMap); err != nil {
			return fmt.Errorf("%w: %w", domain.ErrIntegrationConfigDeserialize, err)
		}
	}

	// Desencriptar credenciales
	var credentials domain.DecryptedCredentials
	usePlatformToken, _ := configMap["use_platform_token"].(bool)

	if usePlatformToken {
		// Obtener credenciales del tipo de integración desde cache
		platformCreds, err := uc.cache.GetPlatformCredentials(ctx, integration.IntegrationTypeID)
		if err != nil || len(platformCreds) == 0 {
			// Fallback: obtener desde DB y desencriptar
			intType, err := uc.repo.GetIntegrationTypeByID(ctx, integration.IntegrationTypeID)
			if err != nil {
				return fmt.Errorf("error al obtener tipo de integración: %w", err)
			}
			if len(intType.PlatformCredentialsEncrypted) > 0 {
				platformCreds, err = uc.encryption.DecryptCredentials(ctx, intType.PlatformCredentialsEncrypted)
				if err != nil {
					return fmt.Errorf("error al desencriptar credenciales de plataforma: %w", err)
				}
			}
		}
		credentials = platformCreds
		uc.log.Info(ctx).
			Uint("integration_type_id", integration.IntegrationTypeID).
			Msg("use_platform_token=true, usando credenciales del tipo de integración")
	} else if len(integration.Credentials) > 0 {
		encryptedBytes, err := decodeEncryptedCredentials([]byte(integration.Credentials))
		if err != nil {
			return fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsDecrypt, err)
		}
		decrypted, err := uc.encryption.DecryptCredentials(ctx, encryptedBytes)
		if err != nil {
			return fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsDecrypt, err)
		}
		credentials = decrypted
	}

	// Obtener el tipo de integración completo (para código y URLs)
	integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, integration.IntegrationTypeID)
	if err != nil {
		return fmt.Errorf("error al obtener tipo de integración: %w", err)
	}
	integrationTypeCode := integrationType.Code

	// Inyectar base_url del tipo si no viene en el config
	if configMap == nil {
		configMap = make(map[string]interface{})
	}
	if _, has := configMap["base_url"]; !has && integrationType.BaseURL != "" {
		configMap["base_url"] = integrationType.BaseURL
	}
	if _, has := configMap["base_url_test"]; !has && integrationType.BaseURLTest != "" {
		configMap["base_url_test"] = integrationType.BaseURLTest
	}
	// Inyectar phone_number_id y whatsapp_url desde config del tipo si use_platform_token
	if usePlatformToken {
		if _, has := configMap["phone_number_id"]; !has {
			var typeConfig map[string]interface{}
			if len(integrationType.ConfigSchema) > 0 {
				json.Unmarshal(integrationType.ConfigSchema, &typeConfig)
			}
			if phoneID, ok := credentials["phone_number_id"]; ok {
				configMap["phone_number_id"] = phoneID
			}
			if waURL, ok := credentials["whatsapp_url"]; ok {
				configMap["whatsapp_url"] = waURL
			}
		}
	}

	// Convertir código string a int
	integrationTypeInt := domain.IntegrationTypeCodeAsInt(integrationTypeCode)

	// Obtener provider registrado para este tipo
	registeredTypes := uc.providerReg.ListRegisteredTypes()
	registeredTypesStr := make([]string, len(registeredTypes))
	for i, t := range registeredTypes {
		registeredTypesStr[i] = fmt.Sprintf("%d", t)
	}
	uc.log.Info(ctx).
		Str("type_code", integrationTypeCode).
		Int("type_int", integrationTypeInt).
		Strs("registered_types", registeredTypesStr).
		Msg("Buscando provider registrado para tipo de integración")

	provider, hasProvider := uc.providerReg.Get(integrationTypeInt)
	if !hasProvider {
		uc.log.Warn(ctx).
			Str("type_code", integrationTypeCode).
			Int("type_int", integrationTypeInt).
			Strs("registered_types", registeredTypesStr).
			Msg("No hay provider registrado, solo validando credenciales básicas")
		return uc.validateBasicCredentials(ctx, integrationTypeCode, credentials)
	}

	uc.log.Info(ctx).
		Str("type_code", integrationTypeCode).
		Msg("Provider encontrado, ejecutando test de conexión real")

	// Para integraciones existentes, verificar que test_phone_number esté presente
	if testPhone, ok := configMap["test_phone_number"].(string); !ok || testPhone == "" {
		uc.log.Error(ctx).
			Uint("id", integration.ID).
			Str("type_code", integrationTypeCode).
			Msg("test_phone_number no encontrado en la configuración de la integración")
		return fmt.Errorf("%w: no hay número de prueba guardado en la configuración de la integración. Por favor, edita la integración y guarda un número de prueba en el campo 'Número de Prueba'", domain.ErrIntegrationTestFailed)
	}

	// Llamar al provider específico
	if err := provider.TestConnection(ctx, configMap, credentials); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Uint("id", integration.ID).
			Str("type_code", integrationTypeCode).
			Msg("Test de conexión falló")
		return fmt.Errorf("%w: %w", domain.ErrIntegrationTestFailed, err)
	}

	uc.log.Info(ctx).
		Uint("id", integration.ID).
		Str("type_code", integrationTypeCode).
		Msg("Test de conexión exitoso")

	return nil
}

// validateBasicCredentials valida credenciales básicas cuando no hay provider registrado
func (uc *IntegrationUseCase) validateBasicCredentials(ctx context.Context, integrationType string, credentials domain.DecryptedCredentials) error {
	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		return domain.ErrIntegrationAccessTokenNotFound
	}

	uc.log.Info(ctx).
		Str("type", integrationType).
		Msg("Validación básica de credenciales exitosa (sin provider registrado)")

	return nil
}

// TestConnectionRaw prueba la conexión con datos proporcionados directamente (sin guardar en BD)
func (uc *IntegrationUseCase) TestConnectionRaw(ctx context.Context, integrationTypeCode string, config map[string]interface{}, credentials map[string]interface{}) error {
	ctx = log.WithFunctionCtx(ctx, "TestConnectionRaw")

	// Convertir código string a int
	integrationTypeInt := domain.IntegrationTypeCodeAsInt(integrationTypeCode)

	// Obtener provider registrado para este tipo
	registeredTypes := uc.providerReg.ListRegisteredTypes()
	registeredTypesStr := make([]string, len(registeredTypes))
	for i, t := range registeredTypes {
		registeredTypesStr[i] = fmt.Sprintf("%d", t)
	}
	uc.log.Info(ctx).
		Str("type_code", integrationTypeCode).
		Int("type_int", integrationTypeInt).
		Strs("registered_types", registeredTypesStr).
		Msg("Buscando provider registrado para tipo de integración (TestConnectionRaw)")

	// Inyectar base_url y base_url_test desde integration_types si el config no los incluye
	// El usuario no ingresa estas URLs — vienen del sistema (integration_types table)
	integrationType, err := uc.repo.GetIntegrationTypeByCode(ctx, integrationTypeCode)
	if err == nil && integrationType != nil {
		if config == nil {
			config = make(map[string]interface{})
		}
		if _, has := config["base_url"]; !has && integrationType.BaseURL != "" {
			config["base_url"] = integrationType.BaseURL
		}
		if _, has := config["base_url_test"]; !has && integrationType.BaseURLTest != "" {
			config["base_url_test"] = integrationType.BaseURLTest
		}
	}

	provider, hasProvider := uc.providerReg.Get(integrationTypeInt)
	if !hasProvider {
		uc.log.Warn(ctx).
			Str("type_code", integrationTypeCode).
			Int("type_int", integrationTypeInt).
			Strs("registered_types", registeredTypesStr).
			Msg("No hay provider registrado, solo validando credenciales básicas")

		// Si usa token de plataforma, no hay credenciales propias que validar
		usePlatformToken, _ := config["use_platform_token"].(bool)
		if usePlatformToken {
			uc.log.Info(ctx).
				Str("type_code", integrationTypeCode).
				Msg("use_platform_token=true, omitiendo validación de credenciales propias")
			return nil
		}
		return uc.validateBasicCredentials(ctx, integrationTypeCode, credentials)
	}

	uc.log.Info(ctx).
		Str("type_code", integrationTypeCode).
		Msg("Provider encontrado, ejecutando test de conexión real")

	// Llamar al provider específico
	if err := provider.TestConnection(ctx, config, credentials); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Str("type_code", integrationTypeCode).
			Msg("Test de conexión falló")
		return fmt.Errorf("%w: %w", domain.ErrIntegrationTestFailed, err)
	}

	uc.log.Info(ctx).
		Str("type_code", integrationTypeCode).
		Msg("Test de conexión exitoso")

	return nil
}
