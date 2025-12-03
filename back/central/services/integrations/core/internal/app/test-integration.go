package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// TestIntegration prueba la conexión de una integración usando el tester registrado
func (uc *integrationUseCase) TestIntegration(ctx context.Context, id uint) error {
	ctx = log.WithFunctionCtx(ctx, "TestIntegration")

	// Obtener integración
	integration, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("integración no encontrada: %w", err)
	}

	// Desencriptar credenciales
	var credentials domain.DecryptedCredentials
	if len(integration.Credentials) > 0 {
		decrypted, err := uc.encryption.DecryptCredentials(ctx, []byte(integration.Credentials))
		if err != nil {
			return fmt.Errorf("error al desencriptar credenciales: %w", err)
		}
		credentials = decrypted
	}

	// Convertir Config a map
	var configMap map[string]interface{}
	if len(integration.Config) > 0 {
		if err := json.Unmarshal(integration.Config, &configMap); err != nil {
			return fmt.Errorf("error al deserializar configuración: %w", err)
		}
	}

	// Obtener tester registrado para este tipo
	tester, err := uc.testerReg.GetTester(integration.Type)
	if err != nil {
		uc.log.Warn(ctx).
			Str("type", integration.Type).
			Msg("No hay tester registrado, solo validando credenciales básicas")
		// Fallback: validación básica si no hay tester
		return uc.validateBasicCredentials(ctx, integration.Type, credentials)
	}

	// Llamar al tester específico
	if err := tester.TestConnection(ctx, configMap, credentials); err != nil {
		uc.log.Error(ctx).
			Err(err).
			Uint("id", integration.ID).
			Str("type", integration.Type).
			Msg("Test de conexión falló")
		return fmt.Errorf("test de conexión falló: %w", err)
	}

	uc.log.Info(ctx).
		Uint("id", integration.ID).
		Str("type", integration.Type).
		Msg("Test de conexión exitoso")

	return nil
}

// validateBasicCredentials valida credenciales básicas cuando no hay tester registrado
func (uc *integrationUseCase) validateBasicCredentials(ctx context.Context, integrationType string, credentials domain.DecryptedCredentials) error {
	// Validación básica: verificar que exista access_token
	accessToken, ok := credentials["access_token"].(string)
	if !ok || accessToken == "" {
		return fmt.Errorf("access_token no encontrado o inválido en las credenciales")
	}

	uc.log.Info(ctx).
		Str("type", integrationType).
		Msg("Validación básica de credenciales exitosa (sin tester registrado)")

	return nil
}
