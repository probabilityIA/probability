package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	boldErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const boldIntegrationTypeCode = "bold_pay"
const boldDefaultBaseURL = "https://integrations.api.bold.co"

type IntegrationRepository struct {
	coreSvc core.IIntegrationCore
	log     log.ILogger
}

func New(coreSvc core.IIntegrationCore, logger log.ILogger) ports.IIntegrationRepository {
	return &IntegrationRepository{
		coreSvc: coreSvc,
		log:     logger.WithModule("bold.integration_repository"),
	}
}

func (r *IntegrationRepository) GetBoldConfig(ctx context.Context) (*ports.BoldConfig, error) {
	return r.GetBoldConfigForMode(ctx, false)
}

func (r *IntegrationRepository) GetBoldConfigForMode(ctx context.Context, testMode bool) (*ports.BoldConfig, error) {
	intType, err := r.coreSvc.GetIntegrationTypeByCode(ctx, boldIntegrationTypeCode)
	if err != nil {
		return nil, boldErrors.ErrBoldConfigNotFound
	}

	creds, err := r.coreSvc.GetCachedPlatformCredentials(ctx, intType.ID)
	if err != nil {
		return nil, fmt.Errorf("get bold platform credentials: %w", err)
	}

	prodAPIKey, _ := creds["api_key"].(string)
	prodSecret, _ := creds["secret_key"].(string)
	testAPIKey, _ := creds["test_api_key"].(string)
	testSecret, _ := creds["test_secret_key"].(string)

	if testMode && testAPIKey == "" && prodAPIKey != "" {
		r.log.Warn(ctx).
			Uint("integration_type_id", intType.ID).
			Msg("Bold test mode requested but test credentials not configured - falling back to production credentials")
		testMode = false
	}
	if !testMode && prodAPIKey == "" && testAPIKey != "" {
		r.log.Warn(ctx).
			Uint("integration_type_id", intType.ID).
			Msg("Bold prod mode requested but prod credentials not configured - using sandbox credentials")
		testMode = true
	}

	var apiKey, secret, environment, baseURL string
	if testMode {
		apiKey = testAPIKey
		secret = testSecret
		environment = "sandbox"
		baseURL = intType.BaseURLTest
		if baseURL == "" {
			baseURL = boldDefaultBaseURL
		}
	} else {
		apiKey = prodAPIKey
		secret = prodSecret
		environment = "production"
		baseURL = intType.BaseURL
		if baseURL == "" {
			baseURL = boldDefaultBaseURL
		}
	}

	if apiKey == "" {
		return nil, boldErrors.ErrInvalidCredentials
	}

	r.log.Info(ctx).
		Uint("integration_type_id", intType.ID).
		Str("environment", environment).
		Str("base_url", baseURL).
		Msg("Bold config retrieved via core")

	return &ports.BoldConfig{
		APIKey:      apiKey,
		Secret:      secret,
		Environment: environment,
		BaseURL:     baseURL,
	}, nil
}
