package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/exito/internal/domain"
)

// TestConnection verifica que las credenciales de Exito sean validas.
// Extrae seller_id (config) y api_key (credentials).
func (uc *exitoUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	sellerID, err := extractString(config, "seller_id")
	if err != nil {
		return domain.ErrMissingSellerID
	}

	apiKey, err := extractString(credentials, "api_key")
	if err != nil {
		return domain.ErrMissingAPIKey
	}

	if err := uc.client.TestConnection(ctx, apiKey, sellerID); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Exito test connection failed")
		return fmt.Errorf("exito: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Msg("Exito test connection successful")
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
