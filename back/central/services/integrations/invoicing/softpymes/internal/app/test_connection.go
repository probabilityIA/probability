package app

import (
	"context"
	"fmt"
)

// TestConnection valida que las credenciales y configuración provistos sean correctos
// contra la API de Softpymes. Implementa ports.IInvoiceUseCase.TestConnection.
func (uc *invoicingUseCase) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	uc.log.Info(ctx).Msg("🧪 Testing connection with Softpymes API")

	apiKey, okKey := credentials["api_key"].(string)
	apiSecret, okSecret := credentials["api_secret"].(string)
	referer, okReferer := config["referer"].(string)

	uc.log.Info(ctx).
		Bool("has_api_key", okKey && apiKey != "").
		Bool("has_api_secret", okSecret && apiSecret != "").
		Bool("has_referer", okReferer && referer != "").
		Msg("📋 Credentials and config validation")

	if !okKey || apiKey == "" {
		return fmt.Errorf("api_key is required in credentials")
	}
	if !okSecret || apiSecret == "" {
		return fmt.Errorf("api_secret is required in credentials")
	}
	if !okReferer || referer == "" {
		return fmt.Errorf("referer is required in config (identificación de instancia del cliente)")
	}

	// Determinar URL efectiva para el test (producción o testing)
	testBaseURL := ""
	if isTesting, ok := config["is_testing"].(bool); ok && isTesting {
		if u, ok := config["base_url_test"].(string); ok && u != "" {
			testBaseURL = u
		}
	}
	if testBaseURL == "" {
		if u, ok := config["base_url"].(string); ok && u != "" {
			testBaseURL = u
		}
	}

	uc.log.Info(ctx).
		Str("test_base_url", testBaseURL).
		Msg("🔌 Calling client.TestAuthentication...")
	if err := uc.client.TestAuthentication(ctx, apiKey, apiSecret, referer, testBaseURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Softpymes connection test failed")
		return fmt.Errorf("failed to connect to Softpymes: %w", err)
	}

	uc.log.Info(ctx).Msg("✅ Softpymes connection test successful")
	return nil
}
