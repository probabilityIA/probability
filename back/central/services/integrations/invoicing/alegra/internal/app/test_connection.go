package app

import (
	"context"
	"fmt"
)

// TestConnection valida que las credenciales y configuración provistos sean correctos
// contra la API de Alegra. Implementa ports.IInvoiceUseCase.TestConnection.
func (uc *invoicingUseCase) TestConnection(ctx context.Context, _ map[string]interface{}, credentials map[string]interface{}) error {
	uc.log.Info(ctx).Msg("🧪 Testing connection with Alegra API")

	email, okEmail := credentials["email"].(string)
	token, okToken := credentials["token"].(string)
	apiURL, _ := credentials["api_url"].(string) // opcional

	uc.log.Info(ctx).
		Bool("has_email", okEmail && email != "").
		Bool("has_token", okToken && token != "").
		Bool("has_api_url", apiURL != "").
		Msg("📋 Alegra credentials validation")

	if !okEmail || email == "" {
		return fmt.Errorf("el campo email es requerido en las credenciales")
	}
	if !okToken || token == "" {
		return fmt.Errorf("el campo token (API key) es requerido en las credenciales")
	}

	uc.log.Info(ctx).Msg("🔌 Calling Alegra client.TestAuthentication...")
	if err := uc.alegraClient.TestAuthentication(ctx, email, token, apiURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Alegra connection test failed")
		return err
	}

	uc.log.Info(ctx).Msg("✅ Alegra connection test successful")
	return nil
}
