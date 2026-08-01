package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/domain"
)

func (uc *tiktokUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	effectiveURL, err := testConnectionBaseURL(config)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Msg("El tipo de integracion TikTok no tiene la URL configurada en base de datos")
		return err
	}

	clientKey, err := extractString(credentials, "client_key")
	if err != nil {
		return domain.ErrMissingClientKey
	}

	clientSecret, err := extractString(credentials, "client_secret")
	if err != nil {
		return domain.ErrMissingClientSecret
	}

	token, err := uc.client.GetClientToken(ctx, effectiveURL, clientKey, clientSecret)
	if err != nil {
		uc.logger.Error(ctx).Err(err).Msg("TikTok test connection failed")
		return err
	}

	if shopName, _ := config["shop_name"].(string); shopName != "" {
		cred := domain.Credential{AccessToken: token, BaseURL: effectiveURL}
		shops, err := uc.client.QueryShopInfo(ctx, cred, shopName, 1)
		if err != nil {
			uc.logger.Error(ctx).Err(err).Str("shop_name", shopName).Msg("TikTok shop query failed")
			return err
		}
		if len(shops) == 0 {
			return domain.ErrShopNotFound
		}
		uc.logger.Info(ctx).
			Str("shop_name", shops[0].ShopName).
			Int64("shop_id", shops[0].ShopID).
			Msg("TikTok test connection successful")
		return nil
	}

	uc.logger.Info(ctx).Msg("TikTok test connection successful")
	return nil
}
