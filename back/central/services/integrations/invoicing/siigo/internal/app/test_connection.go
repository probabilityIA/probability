package app

import (
	"context"
	"fmt"
)

// TestConnection valida que las credenciales y configuración provistos sean correctos
// contra la API de Siigo. Implementa ports.IInvoiceUseCase.TestConnection.
func (uc *invoicingUseCase) TestConnection(ctx context.Context, _ map[string]interface{}, credentials map[string]interface{}) error {
	uc.log.Info(ctx).Msg("🧪 Testing connection with Siigo API")

	username, okUsername := credentials["username"].(string)
	accessKey, okAccessKey := credentials["access_key"].(string)
	accountID, okAccountID := credentials["account_id"].(string)
	partnerID, okPartnerID := credentials["partner_id"].(string)
	apiURL, _ := credentials["api_url"].(string) // opcional

	uc.log.Info(ctx).
		Bool("has_username", okUsername && username != "").
		Bool("has_access_key", okAccessKey && accessKey != "").
		Bool("has_account_id", okAccountID && accountID != "").
		Bool("has_partner_id", okPartnerID && partnerID != "").
		Bool("has_api_url", apiURL != "").
		Msg("📋 Siigo credentials validation")

	if !okUsername || username == "" {
		return fmt.Errorf("el campo username (email) es requerido en las credenciales")
	}
	if !okAccessKey || accessKey == "" {
		return fmt.Errorf("el campo access_key es requerido en las credenciales")
	}
	if !okAccountID || accountID == "" {
		return fmt.Errorf("el campo account_id (subscription key) es requerido en las credenciales")
	}
	if !okPartnerID || partnerID == "" {
		return fmt.Errorf("el campo partner_id es requerido en las credenciales")
	}

	uc.log.Info(ctx).Msg("🔌 Calling Siigo client.TestAuthentication...")
	if err := uc.siigoClient.TestAuthentication(ctx, username, accessKey, accountID, partnerID, apiURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Siigo connection test failed")
		return err
	}

	uc.log.Info(ctx).Msg("✅ Siigo connection test successful")
	return nil
}
