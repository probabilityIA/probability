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

	// Desencriptar credenciales
	var credentials domain.DecryptedCredentials
	if len(integration.Credentials) > 0 {
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

	// Convertir Config a map
	var configMap map[string]interface{}
	if len(integration.Config) > 0 {
		if err := json.Unmarshal(integration.Config, &configMap); err != nil {
			return fmt.Errorf("%w: %w", domain.ErrIntegrationConfigDeserialize, err)
		}
	}

	// Obtener el código del tipo de integración
	integrationTypeCode := ""
	if integration.IntegrationType != nil {
		integrationTypeCode = integration.IntegrationType.Code
	} else {
		integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, integration.IntegrationTypeID)
		if err != nil {
			return fmt.Errorf("error al obtener tipo de integración: %w", err)
		}
		integrationTypeCode = integrationType.Code
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
