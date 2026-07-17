package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
)

func (uc *jumpsellerUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	effectiveURL, err := testConnectionBaseURL(config)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Msg("El tipo de integracion Jumpseller no tiene la URL configurada en base de datos")
		return err
	}

	cred := domain.Credential{BaseURL: effectiveURL}

	if accessToken, _ := credentials["access_token"].(string); accessToken != "" {
		cred.AccessToken = accessToken
	} else {
		apiKey, err := extractString(credentials, "api_key")
		if err != nil {
			return domain.ErrMissingAPIKey
		}
		apiSecret, err := extractString(credentials, "api_secret")
		if err != nil {
			return domain.ErrMissingAPISecret
		}
		cred.APIKey = apiKey
		cred.APISecret = apiSecret
	}

	storeInfo, err := uc.client.GetStoreInfo(ctx, cred)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Msg("Jumpseller test connection failed")
		return err
	}

	uc.logger.Info(ctx).
		Str("store_code", storeInfo.Code).
		Str("store_name", storeInfo.Name).
		Msg("Jumpseller test connection successful")

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
