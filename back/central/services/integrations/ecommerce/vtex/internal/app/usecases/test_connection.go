package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

// TestConnection verifica que las credenciales de VTEX sean válidas.
// Extrae store_url (config) y api_key / api_token (credentials).
func (uc *vtexUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	storeURL, err := extractString(config, "store_url")
	if err != nil {
		return domain.ErrMissingStoreURL
	}

	apiKey, err := extractString(credentials, "api_key")
	if err != nil {
		return domain.ErrMissingAPIKey
	}

	apiToken, err := extractString(credentials, "api_token")
	if err != nil {
		return domain.ErrMissingAPIToken
	}

	if err := uc.client.TestConnection(ctx, storeURL, apiKey, apiToken); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("VTEX test connection failed")
		return fmt.Errorf("vtex: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Msg("VTEX test connection successful")
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
