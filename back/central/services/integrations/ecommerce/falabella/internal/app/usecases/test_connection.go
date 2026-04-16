package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/falabella/internal/domain"
)

// TestConnection verifica que las credenciales de Falabella Seller Center sean válidas.
// Extrae user_id (config) y api_key (credentials).
func (uc *falabellaUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	userID, err := extractString(config, "user_id")
	if err != nil {
		return domain.ErrMissingUserID
	}

	apiKey, err := extractString(credentials, "api_key")
	if err != nil {
		return domain.ErrMissingAPIKey
	}

	if err := uc.client.TestConnection(ctx, apiKey, userID); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Falabella test connection failed")
		return fmt.Errorf("falabella: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Msg("Falabella test connection successful")
	return nil
}

// extractString extrae un campo string de un mapa, retornando error si falta o es vacío.
func extractString(m map[string]interface{}, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("missing field: %s", key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("field %s must be a non-empty string", key)
	}
	return s, nil
}
