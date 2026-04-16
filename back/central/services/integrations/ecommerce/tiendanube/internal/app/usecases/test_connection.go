package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

// TestConnection verifica que las credenciales de Tiendanube sean validas.
// Extrae store_url (config) y access_token (credentials).
func (uc *tiendanubeUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	storeURL, err := extractString(config, "store_url")
	if err != nil {
		return domain.ErrMissingStoreURL
	}

	accessToken, err := extractString(credentials, "access_token")
	if err != nil {
		return domain.ErrMissingAccessToken
	}

	if err := uc.client.TestConnection(ctx, storeURL, accessToken); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Tiendanube test connection failed")
		return fmt.Errorf("tiendanube: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Msg("Tiendanube test connection successful")
	return nil
}

// extractString extrae un campo string de un mapa, retornando error si falta o es vacio.
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
