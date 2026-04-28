package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/errors"
)

const boldIntegrationTypeCode = "bold_pay"

func (r *Repository) GetBoldCredentials(ctx context.Context) (*dtos.BoldCredentials, error) {
	return r.getBoldCredentials(ctx, nil)
}

func (r *Repository) GetBoldCredentialsForBusiness(ctx context.Context, businessID uint) (*dtos.BoldCredentials, error) {
	biz, err := r.GetBoldIntegrationForBusiness(ctx, businessID)
	if err != nil {
		return nil, err
	}
	isTesting := biz != nil && biz.IsTesting
	return r.getBoldCredentials(ctx, &isTesting)
}

func (r *Repository) getBoldCredentials(ctx context.Context, forceTesting *bool) (*dtos.BoldCredentials, error) {
	if r.integrationCore == nil {
		return nil, domainerrors.ErrBoldConfigNotFound
	}

	intType, err := r.integrationCore.GetIntegrationTypeByCode(ctx, boldIntegrationTypeCode)
	if err != nil {
		return nil, domainerrors.ErrBoldConfigNotFound
	}

	creds, err := r.integrationCore.GetCachedPlatformCredentials(ctx, intType.ID)
	if err != nil {
		return nil, domainerrors.ErrBoldCredentialsMissing
	}

	prodAPIKey, _ := creds["api_key"].(string)
	prodSecretKey, _ := creds["secret_key"].(string)
	testAPIKey, _ := creds["test_api_key"].(string)
	testSecretKey, _ := creds["test_secret_key"].(string)

	var environment string
	switch {
	case forceTesting != nil && *forceTesting:
		environment = "sandbox"
	case forceTesting != nil && !*forceTesting:
		environment = "production"
	default:
		environment, _ = creds["environment"].(string)
		if environment == "" {
			if testAPIKey != "" && testSecretKey != "" {
				environment = "sandbox"
			} else {
				environment = "production"
			}
		}
	}

	apiKey, secretKey := prodAPIKey, prodSecretKey
	if environment == "sandbox" && testAPIKey != "" && testSecretKey != "" {
		apiKey, secretKey = testAPIKey, testSecretKey
	}
	if apiKey == "" || secretKey == "" {
		return nil, domainerrors.ErrBoldCredentialsMissing
	}

	baseURL := intType.BaseURL
	if environment == "sandbox" && intType.BaseURLTest != "" {
		baseURL = intType.BaseURLTest
	}
	if baseURL == "" {
		baseURL = "https://integrations.api.bold.co"
	}

	r.log.Info(ctx).
		Uint("integration_type_id", intType.ID).
		Str("environment", environment).
		Msg("Bold credentials retrieved via core")

	return &dtos.BoldCredentials{
		APIKey:            apiKey,
		SecretKey:         secretKey,
		Environment:       environment,
		BaseURL:           baseURL,
		IntegrationTypeID: intType.ID,
	}, nil
}

func (r *Repository) GetBoldIntegrationForBusiness(ctx context.Context, businessID uint) (*dtos.BoldBusinessIntegration, error) {
	if r.integrationCore == nil {
		return nil, domainerrors.ErrBoldConfigNotFound
	}
	intType, err := r.integrationCore.GetIntegrationTypeByCode(ctx, boldIntegrationTypeCode)
	if err != nil {
		return nil, domainerrors.ErrBoldConfigNotFound
	}

	var row struct {
		ID        uint
		IsTesting bool
	}
	err = r.db.Conn(ctx).Table("integrations").
		Select("id, is_testing").
		Where("business_id = ? AND integration_type_id = ? AND is_active = true", businessID, intType.ID).
		Limit(1).
		Take(&row).Error
	if err != nil {
		return &dtos.BoldBusinessIntegration{IntegrationTypeID: intType.ID}, nil
	}
	return &dtos.BoldBusinessIntegration{
		IntegrationID:     row.ID,
		IntegrationTypeID: intType.ID,
		IsTesting:         row.IsTesting,
	}, nil
}
