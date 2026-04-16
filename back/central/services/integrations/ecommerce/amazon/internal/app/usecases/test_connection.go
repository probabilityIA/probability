package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/amazon/internal/domain"
)

// TestConnection verifica que las credenciales de Amazon SP-API sean validas.
// Extrae seller_id (config) y refresh_token, client_id, client_secret (credentials).
func (uc *amazonUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	sellerID, err := extractString(config, "seller_id")
	if err != nil {
		return domain.ErrMissingSellerID
	}

	refreshToken, err := extractString(credentials, "refresh_token")
	if err != nil {
		return domain.ErrMissingRefreshToken
	}

	clientID, err := extractString(credentials, "client_id")
	if err != nil {
		return domain.ErrMissingClientID
	}

	clientSecret, err := extractString(credentials, "client_secret")
	if err != nil {
		return domain.ErrMissingClientSecret
	}

	if err := uc.client.TestConnection(ctx, sellerID, refreshToken, clientID, clientSecret); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Amazon test connection failed")
		return fmt.Errorf("amazon: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Msg("Amazon test connection successful")
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
