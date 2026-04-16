package app

import (
	"context"
	"fmt"
)

// TestConnection valida que las credenciales y configuración provistos sean correctos
// contra la API de Factus. Implementa ports.IInvoiceUseCase.TestConnection.
func (uc *invoicingUseCase) TestConnection(ctx context.Context, _ map[string]interface{}, credentials map[string]interface{}) error {
	uc.log.Info(ctx).Msg("🧪 Testing connection with Factus API")

	clientID, okClientID := credentials["client_id"].(string)
	clientSecret, okClientSecret := credentials["client_secret"].(string)
	username, okUsername := credentials["username"].(string)
	password, okPassword := credentials["password"].(string)
	apiURL, _ := credentials["api_url"].(string) // opcional

	uc.log.Info(ctx).
		Bool("has_client_id", okClientID && clientID != "").
		Bool("has_client_secret", okClientSecret && clientSecret != "").
		Bool("has_username", okUsername && username != "").
		Bool("has_password", okPassword && password != "").
		Bool("has_api_url", apiURL != "").
		Msg("📋 Factus credentials validation")

	if !okClientID || clientID == "" {
		return fmt.Errorf("el campo client_id es requerido en las credenciales")
	}
	if !okClientSecret || clientSecret == "" {
		return fmt.Errorf("el campo client_secret es requerido en las credenciales")
	}
	if !okUsername || username == "" {
		return fmt.Errorf("el campo username (email) es requerido en las credenciales")
	}
	if !okPassword || password == "" {
		return fmt.Errorf("el campo password es requerido en las credenciales")
	}

	uc.log.Info(ctx).Msg("🔌 Calling Factus client.TestAuthentication...")
	if err := uc.factusClient.TestAuthentication(ctx, apiURL, clientID, clientSecret, username, password); err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Factus connection test failed")
		return err
	}

	uc.log.Info(ctx).Msg("✅ Factus connection test successful")
	return nil
}
