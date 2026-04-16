package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// TestConnection verifica que las credenciales de WooCommerce sean válidas.
// Extrae store_url (config) y consumer_key / consumer_secret (credentials).
func (uc *wooCommerceUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	storeURL, err := extractString(config, "store_url")
	if err != nil {
		return domain.ErrMissingStoreURL
	}

	consumerKey, err := extractString(credentials, "consumer_key")
	if err != nil {
		return domain.ErrMissingConsumerKey
	}

	consumerSecret, err := extractString(credentials, "consumer_secret")
	if err != nil {
		return domain.ErrMissingConsumerSecret
	}

	if err := uc.client.TestConnection(ctx, storeURL, consumerKey, consumerSecret); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("WooCommerce test connection failed")
		return fmt.Errorf("woocommerce: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Msg("WooCommerce test connection successful")
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
