package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// GetPublicIntegrationByID obtiene una integración por su ID en formato público
func (uc *IntegrationUseCase) GetPublicIntegrationByID(ctx context.Context, integrationID string) (*domain.PublicIntegration, error) {
	ctx = log.WithFunctionCtx(ctx, "GetPublicIntegrationByID")

	// Parsear ID de string a uint
	var id uint
	if _, err := fmt.Sscanf(integrationID, "%d", &id); err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Invalid integration ID format")
		return nil, fmt.Errorf("invalid integration ID: %w", err)
	}

	// Obtener integración del dominio
	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener integración")
		return nil, err
	}

	// Mapear a formato público
	return uc.mapToPublicIntegration(integration), nil
}

// GetIntegrationConfig obtiene solo la configuración de una integración por tipo
func (uc *IntegrationUseCase) GetIntegrationConfig(ctx context.Context, integrationType string, businessID *uint) (map[string]interface{}, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationConfig")

	integration, err := uc.GetIntegrationByType(ctx, integrationType, businessID)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if len(integration.Config) > 0 {
		if err := json.Unmarshal(integration.Config, &config); err != nil {
			uc.log.Error(ctx).Err(err).Msg("Error al deserializar configuración")
			return nil, fmt.Errorf("error al deserializar configuración: %w", err)
		}
	}

	return config, nil
}

// DecryptCredentialField desencripta un campo específico de las credenciales
func (uc *IntegrationUseCase) DecryptCredentialField(ctx context.Context, integrationID string, fieldName string) (string, error) {
	ctx = log.WithFunctionCtx(ctx, "DecryptCredentialField")

	// Parsear ID de string a uint
	var id uint
	if _, err := fmt.Sscanf(integrationID, "%d", &id); err != nil {
		uc.log.Error(ctx).Err(err).Str("integration_id", integrationID).Msg("Invalid integration ID format")
		return "", fmt.Errorf("invalid integration ID: %w", err)
	}

	// Obtener integración con credenciales
	integration, err := uc.GetIntegrationByIDWithCredentials(ctx, id)
	if err != nil {
		uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener integración con credenciales")
		return "", err
	}

	// Validar que existan credenciales
	if integration.DecryptedCredentials == nil || len(integration.DecryptedCredentials) == 0 {
		uc.log.Error(ctx).Uint("id", id).Msg("No credentials found for integration")
		return "", fmt.Errorf("no credentials found for integration %d. Please update the integration with credentials", id)
	}

	// Log de debug: mostrar qué campos están disponibles
	availableFields := make([]string, 0, len(integration.DecryptedCredentials))
	for k := range integration.DecryptedCredentials {
		availableFields = append(availableFields, k)
	}
	uc.log.Debug(ctx).
		Str("field", fieldName).
		Uint("id", id).
		Strs("available_fields", availableFields).
		Msg("Attempting to get credential field")

	// Obtener el campo
	value, ok := integration.DecryptedCredentials[fieldName]
	if !ok {
		uc.log.Error(ctx).
			Str("field", fieldName).
			Uint("id", id).
			Strs("available_fields", availableFields).
			Msg("Field not found in credentials")
		if len(availableFields) > 0 {
			return "", fmt.Errorf("field '%s' not found in credentials for integration %d. Available fields: %v. Please update the integration with the correct credentials", fieldName, id, availableFields)
		}
		return "", fmt.Errorf("field '%s' not found in credentials for integration %d. Credentials are empty. Please update the integration with credentials", fieldName, id)
	}

	// Validar que sea string
	strValue, ok := value.(string)
	if !ok {
		uc.log.Error(ctx).Str("field", fieldName).Uint("id", id).Msg("Field is not a string")
		return "", fmt.Errorf("field %s is not a string", fieldName)
	}

	return strValue, nil
}

// GetPlatformCredentialByIntegrationID obtiene un campo específico de las credenciales
// de plataforma del tipo de integración asociado a la integración dada.
// Se usa cuando una integración tiene use_platform_token=true en su configuración.
func (uc *IntegrationUseCase) GetPlatformCredentialByIntegrationID(ctx context.Context, integrationID string, fieldName string) (string, error) {
	ctx = log.WithFunctionCtx(ctx, "GetPlatformCredentialByIntegrationID")

	// Parsear ID
	var id uint
	if _, err := fmt.Sscanf(integrationID, "%d", &id); err != nil {
		return "", fmt.Errorf("invalid integration ID: %w", err)
	}

	// Obtener integración para conocer el IntegrationTypeID
	integration, err := uc.repo.GetIntegrationByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("error al obtener integración %s: %w", integrationID, err)
	}

	// Obtener tipo de integración (contiene las credenciales de plataforma encriptadas)
	integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, integration.IntegrationTypeID)
	if err != nil {
		return "", fmt.Errorf("error al obtener tipo de integración %d: %w", integration.IntegrationTypeID, err)
	}

	if len(integrationType.PlatformCredentialsEncrypted) == 0 {
		return "", fmt.Errorf("no hay credenciales de plataforma configuradas para el tipo de integración %d", integration.IntegrationTypeID)
	}

	// Desencriptar
	credentials, err := uc.encryption.DecryptCredentials(ctx, integrationType.PlatformCredentialsEncrypted)
	if err != nil {
		return "", fmt.Errorf("error al desencriptar credenciales de plataforma: %w", err)
	}

	value, ok := credentials[fieldName]
	if !ok {
		return "", fmt.Errorf("campo '%s' no encontrado en las credenciales de plataforma del tipo de integración %d", fieldName, integration.IntegrationTypeID)
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("el campo '%s' no es un string en las credenciales de plataforma", fieldName)
	}

	return strValue, nil
}

// mapToPublicIntegration mapea una integración del dominio a formato público
func (uc *IntegrationUseCase) mapToPublicIntegration(integration *domain.Integration) *domain.PublicIntegration {
	var config map[string]interface{}
	if len(integration.Config) > 0 {
		_ = json.Unmarshal(integration.Config, &config)
	}

	// Usar directamente el IntegrationTypeID (es el ID de la tabla integration_types)
	integrationTypeID := int(integration.IntegrationTypeID)

	pub := &domain.PublicIntegration{
		ID:              integration.ID,
		BusinessID:      integration.BusinessID,
		Name:            integration.Name,
		StoreID:         integration.StoreID,
		IntegrationType: integrationTypeID,
		Config:          config,
		IsTesting:       integration.IsTesting,
	}

	if integration.IntegrationType != nil {
		pub.BaseURL = integration.IntegrationType.BaseURL
		pub.BaseURLTest = integration.IntegrationType.BaseURLTest
	}

	return pub
}
