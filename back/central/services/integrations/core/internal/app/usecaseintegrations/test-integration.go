package usecaseintegrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// TestIntegration prueba la conexión de una integración usando el tester registrado
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
		decrypted, err := uc.encryption.DecryptCredentials(ctx, []byte(integration.Credentials))
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

	// Obtener el código del tipo de integración para el tester
	integrationTypeCode := ""
	if integration.IntegrationType != nil {
		integrationTypeCode = integration.IntegrationType.Code
	} else {
		// Si no está cargado, obtenerlo del repositorio
		integrationType, err := uc.repo.GetIntegrationTypeByID(ctx, integration.IntegrationTypeID)
		if err != nil {
			return fmt.Errorf("error al obtener tipo de integración: %w", err)
		}
		integrationTypeCode = integrationType.Code
	}

	// Obtener tester registrado para este tipo
	tester, err := uc.testerReg.GetTester(integrationTypeCode)
	if err != nil {
		uc.log.Warn(ctx).
			Str("type_code", integrationTypeCode).
			Msg("No hay tester registrado, solo validando credenciales básicas")
		// Fallback: validación básica si no hay tester
		return uc.validateBasicCredentials(ctx, integrationTypeCode, credentials)
	}

	// Llamar al tester específico
	if err := tester.TestConnection(ctx, configMap, credentials); err != nil {
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

// validateBasicCredentials valida credenciales básicas cuando no hay tester registrado
func (uc *IntegrationUseCase) validateBasicCredentials(ctx context.Context, integrationType string, credentials domain.DecryptedCredentials) error {
	// Validación básica: verificar que exista access_token
	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		return domain.ErrIntegrationAccessTokenNotFound
	}

	uc.log.Info(ctx).
		Str("type", integrationType).
		Msg("Validación básica de credenciales exitosa (sin tester registrado)")

	return nil
}
