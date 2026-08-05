package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	bcErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	integrationTypeCode   = "bancolombia_qr"
	defaultSandboxBaseURL = "https://sb-api.bancolombia.com"
	defaultProdBaseURL    = "https://api.bancolombia.com"
)

type IntegrationRepository struct {
	coreSvc core.IIntegrationCore
	log     log.ILogger
}

func New(coreSvc core.IIntegrationCore, logger log.ILogger) ports.IIntegrationRepository {
	return &IntegrationRepository{
		coreSvc: coreSvc,
		log:     logger.WithModule("bancolombia.integration_repository"),
	}
}

func (r *IntegrationRepository) GetConfig(ctx context.Context) (*ports.BancolombiaConfig, error) {
	return r.GetConfigForMode(ctx, false)
}

func (r *IntegrationRepository) GetConfigForMode(ctx context.Context, testMode bool) (*ports.BancolombiaConfig, error) {
	intType, err := r.coreSvc.GetIntegrationTypeByCode(ctx, integrationTypeCode)
	if err != nil {
		return nil, bcErrors.ErrConfigNotFound
	}

	creds, err := r.coreSvc.GetCachedPlatformCredentials(ctx, intType.ID)
	if err != nil {
		return nil, fmt.Errorf("get bancolombia platform credentials: %w", err)
	}

	str := func(key string) string {
		v, _ := creds[key].(string)
		return v
	}

	prodClientID := str("client_id")
	prodClientSecret := str("client_secret")
	testClientID := str("test_client_id")
	testClientSecret := str("test_client_secret")

	if testMode && testClientID == "" && prodClientID != "" {
		r.log.Warn(ctx).
			Uint("integration_type_id", intType.ID).
			Msg("Bancolombia test mode requested but test credentials not configured - falling back to production credentials")
		testMode = false
	}
	if !testMode && prodClientID == "" && testClientID != "" {
		r.log.Warn(ctx).
			Uint("integration_type_id", intType.ID).
			Msg("Bancolombia prod mode requested but prod credentials not configured - using sandbox credentials")
		testMode = true
	}

	var config *ports.BancolombiaConfig
	if testMode {
		baseURL := intType.BaseURLTest
		if baseURL == "" {
			baseURL = defaultSandboxBaseURL
		}
		config = &ports.BancolombiaConfig{
			ClientID:      testClientID,
			ClientSecret:  testClientSecret,
			WebhookSecret: str("test_webhook_secret"),
			Environment:   "sandbox",
			BaseURL:       baseURL,
		}
	} else {
		baseURL := intType.BaseURL
		if baseURL == "" {
			baseURL = defaultProdBaseURL
		}
		config = &ports.BancolombiaConfig{
			ClientID:      prodClientID,
			ClientSecret:  prodClientSecret,
			WebhookSecret: str("webhook_secret"),
			Environment:   "production",
			BaseURL:       baseURL,
		}
	}

	config.TokenPath = str("token_path")
	config.QRGeneratePath = str("qr_generate_path")
	config.QRQueryPath = str("qr_query_path")
	config.Scope = str("scope")
	config.MerchantID = str("merchant_id")
	config.MerchantName = str("merchant_name")

	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, bcErrors.ErrInvalidCredentials
	}

	r.log.Info(ctx).
		Uint("integration_type_id", intType.ID).
		Str("environment", config.Environment).
		Str("base_url", config.BaseURL).
		Msg("Bancolombia config retrieved via core")

	return config, nil
}
