package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func (uc *vtexUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	accountName, err := extractString(config, "account_name")
	if err != nil {
		return domain.ErrMissingAccountName
	}

	accountName = CleanAccountName(accountName)
	if accountName == "" {
		return domain.ErrMissingAccountName
	}

	appKey, err := extractString(credentials, "app_key")
	if err != nil {
		return domain.ErrMissingAppKey
	}

	appToken, err := extractString(credentials, "app_token")
	if err != nil {
		return domain.ErrMissingAppToken
	}

	cred := domain.Credential{
		AccountName: accountName,
		AppKey:      appKey,
		AppToken:    appToken,
	}

	if err := uc.client.TestConnection(ctx, cred); err != nil {
		uc.logger.Error(ctx).Err(err).Str("account", accountName).Msg("VTEX test connection failed")
		return fmt.Errorf("vtex: test connection failed: %w", err)
	}

	uc.logger.Info(ctx).Str("account", accountName).Msg("VTEX test connection successful")
	return nil
}
