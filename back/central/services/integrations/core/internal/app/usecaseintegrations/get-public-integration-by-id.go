package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// PublicIntegration representa una integración en formato público
type PublicIntegration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
}

// GetPublicIntegrationByID obtiene una integración por su ID en formato público
func (uc *IntegrationUseCase) GetPublicIntegrationByID(ctx context.Context, integrationID string) (*PublicIntegration, error) {
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
	if integration.DecryptedCredentials == nil {
		uc.log.Error(ctx).Uint("id", id).Msg("No credentials found for integration")
		return "", fmt.Errorf("no credentials found for integration")
	}

	// Obtener el campo
	value, ok := integration.DecryptedCredentials[fieldName]
	if !ok {
		uc.log.Error(ctx).Str("field", fieldName).Uint("id", id).Msg("Field not found in credentials")
		return "", fmt.Errorf("field %s not found in credentials", fieldName)
	}

	// Validar que sea string
	strValue, ok := value.(string)
	if !ok {
		uc.log.Error(ctx).Str("field", fieldName).Uint("id", id).Msg("Field is not a string")
		return "", fmt.Errorf("field %s is not a string", fieldName)
	}

	return strValue, nil
}

// MapToPublicIntegration mapea una integración del dominio a formato público (método público para uso interno)
func (uc *IntegrationUseCase) MapToPublicIntegration(integration *domain.Integration) *PublicIntegration {
	return uc.mapToPublicIntegration(integration)
}

// mapToPublicIntegration mapea una integración del dominio a formato público
func (uc *IntegrationUseCase) mapToPublicIntegration(integration *domain.Integration) *PublicIntegration {
	var config map[string]interface{}
	if len(integration.Config) > 0 {
		_ = json.Unmarshal(integration.Config, &config)
	}

	// Obtener el código del tipo de integración y convertirlo a int
	integrationTypeCode := 0
	if integration.IntegrationType != nil {
		integrationTypeCode = getIntegrationTypeCodeAsInt(integration.IntegrationType.Code)
	}

	return &PublicIntegration{
		ID:              integration.ID,
		BusinessID:      integration.BusinessID,
		Name:            integration.Name,
		StoreID:         integration.StoreID,
		IntegrationType: integrationTypeCode,
		Config:          config,
	}
}

// getIntegrationTypeCodeAsInt convierte el código de tipo de integración a int
// Esta función mapea los códigos antiguos (strings) a los nuevos (int)
func getIntegrationTypeCodeAsInt(code string) int {
	switch code {
	case "shopify":
		return 1 // IntegrationTypeShopify
	case "whatsapp":
		return 2 // IntegrationTypeWhatsApp
	case "mercado_libre":
		return 3 // IntegrationTypeMercadoLibre
	default:
		return 0
	}
}
