package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func (uc *meliUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	accessToken, err := extractString(credentials, "access_token")
	if err != nil {
		return domain.ErrMissingAccessToken
	}

	cli := uc.client
	if isTesting, _ := config["is_testing"].(bool); isTesting {
		if raw, terr := extractString(config, "base_url_test"); terr == nil {
			if baseURL := normalizeBaseURL(raw); baseURL != "" {
				uc.logger.Info(ctx).Str("base_url", baseURL).Msg("Meli test connection en modo pruebas")
				cli = uc.client.WithBaseURL(baseURL)
			}
		}
	}

	if err := cli.TestConnection(ctx, accessToken); err != nil {
		uc.logger.Error(ctx).Err(err).Msg("MercadoLibre test connection failed")
		return fmt.Errorf("meli: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Msg("MercadoLibre test connection successful")
	return nil
}

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
