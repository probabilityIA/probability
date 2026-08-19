package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func (uc *tiendanubeUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	baseURL, err := testConnectionBaseURL(config)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Msg("El tipo de integracion Tiendanube no tiene la URL configurada en base de datos")
		return err
	}

	accessToken, err := extractString(credentials, "access_token")
	if err != nil {
		return domain.ErrMissingAccessToken
	}

	storeID := resolveStoreID(config, "")
	if storeID == "" {
		return domain.ErrMissingStoreID
	}

	cred := domain.Credential{
		AccessToken: accessToken,
		StoreID:     storeID,
		BaseURL:     baseURL,
		UserAgent:   optionalString(config, configUserAgent),
	}

	storeInfo, err := uc.client.GetStoreInfo(ctx, cred)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Tiendanube test connection failed")
		return err
	}

	uc.logger.Info(ctx).
		Str("store_name", storeInfo.Name).
		Str("store_url", storeInfo.URL).
		Msg("Tiendanube test connection successful")

	return nil
}

func testConnectionBaseURL(config map[string]interface{}) (string, error) {
	if isTesting, _ := config["is_testing"].(bool); isTesting {
		url, err := extractString(config, "base_url_test")
		if err != nil {
			return "", domain.ErrMissingBaseURLTest
		}
		return url, nil
	}
	url, err := extractString(config, "base_url")
	if err != nil {
		return "", domain.ErrMissingBaseURL
	}
	return url, nil
}
