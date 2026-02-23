package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

// TestConnection verifica que las credenciales de MercadoLibre sean válidas.
// Extrae el access_token de credentials y hace una llamada de prueba a la API.
func (uc *meliUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	accessToken, err := extractString(credentials, "access_token")
	if err != nil {
		return domain.ErrMissingAccessToken
	}

	if err := uc.client.TestConnection(ctx, accessToken); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("MercadoLibre test connection failed")
		return fmt.Errorf("meli: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Msg("MercadoLibre test connection successful")
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
