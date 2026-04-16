package usecaseintegrations

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// decodeEncryptedCredentials decodifica las credenciales desde base64
// Las credenciales se guardan como JSON: {"encrypted": "base64string"}
func decodeEncryptedCredentials(encodedJSON []byte) ([]byte, error) {
	var wrapper map[string]string
	if err := json.Unmarshal(encodedJSON, &wrapper); err != nil {
		// Si no es JSON válido, asumir que es el formato antiguo (bytes directos)
		// Intentar decodificar como base64 directamente
		if decoded, err := base64.StdEncoding.DecodeString(string(encodedJSON)); err == nil {
			return decoded, nil
		}
		return encodedJSON, nil // Retornar como está si no se puede decodificar
	}

	if encrypted, ok := wrapper["encrypted"]; ok {
		decoded, err := base64.StdEncoding.DecodeString(encrypted)
		if err != nil {
			return nil, fmt.Errorf("error al decodificar base64: %w", err)
		}
		return decoded, nil
	}

	return nil, fmt.Errorf("campo 'encrypted' no encontrado en credenciales")
}

// GetIntegrationByType obtiene una integración por código de tipo y business_id, con credenciales desencriptadas
func (uc *IntegrationUseCase) GetIntegrationByType(ctx context.Context, integrationTypeCode string, businessID *uint) (*domain.IntegrationWithCredentials, error) {
	ctx = log.WithFunctionCtx(ctx, "GetIntegrationByType")

	// Primero obtener el IntegrationType por código
	integrationType, err := uc.repo.GetIntegrationTypeByCode(ctx, integrationTypeCode)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Str("type_code", integrationTypeCode).
			Msg("Error al obtener tipo de integración por código")
		return nil, fmt.Errorf("%w '%s': %w", domain.ErrIntegrationTypeNotFound, integrationTypeCode, err)
	}

	// Obtener integración del repository usando el IntegrationTypeID
	integration, err := uc.repo.GetActiveIntegrationByIntegrationTypeID(ctx, integrationType.ID, businessID)
	if err != nil {
		uc.log.Error(ctx).Err(err).
			Uint("integration_type_id", integrationType.ID).
			Str("type_code", integrationTypeCode).
			Msg("Error al obtener integración por tipo")
		return nil, err
	}

	// Desencriptar credenciales
	var decryptedCredentials domain.DecryptedCredentials
	if len(integration.Credentials) > 0 {
		// Las credenciales están codificadas en base64 dentro de un JSON
		encryptedBytes, err := decodeEncryptedCredentials([]byte(integration.Credentials))
		if err != nil {
			uc.log.Error(ctx).Err(err).
				Uint("id", integration.ID).
				Msg("Error al decodificar credenciales desde base64")
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsDecrypt, err)
		}
		decrypted, err := uc.encryption.DecryptCredentials(ctx, encryptedBytes)
		if err != nil {
			uc.log.Error(ctx).Err(err).
				Uint("id", integration.ID).
				Msg("Error al desencriptar credenciales")
			return nil, fmt.Errorf("%w: %w", domain.ErrIntegrationCredentialsDecrypt, err)
		}
		decryptedCredentials = decrypted
	}

	return &domain.IntegrationWithCredentials{
		Integration:          *integration,
		DecryptedCredentials: decryptedCredentials,
	}, nil
}
