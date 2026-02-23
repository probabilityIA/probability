package app

import (
	"context"
	"fmt"
)

// TestConnection valida que las credenciales y configuración provistos sean correctos
// contra la API de Helisa. Implementa ports.IInvoiceUseCase.TestConnection.
func (uc *invoicingUseCase) TestConnection(ctx context.Context, _ map[string]interface{}, credentials map[string]interface{}) error {
	uc.log.Info(ctx).Msg("🧪 Testing connection with Helisa API")

	username, okUsername := credentials["username"].(string)
	password, okPassword := credentials["password"].(string)
	companyID, okCompanyID := credentials["company_id"].(string)
	apiURL, _ := credentials["api_url"].(string) // opcional

	uc.log.Info(ctx).
		Bool("has_username", okUsername && username != "").
		Bool("has_password", okPassword && password != "").
		Bool("has_company_id", okCompanyID && companyID != "").
		Bool("has_api_url", apiURL != "").
		Msg("📋 Helisa credentials validation")

	if !okUsername || username == "" {
		return fmt.Errorf("el campo username es requerido en las credenciales")
	}
	if !okPassword || password == "" {
		return fmt.Errorf("el campo password es requerido en las credenciales")
	}
	if !okCompanyID || companyID == "" {
		return fmt.Errorf("el campo company_id es requerido en las credenciales")
	}

	uc.log.Info(ctx).Msg("🔌 Calling Helisa client.TestAuthentication...")
	if err := uc.helisaClient.TestAuthentication(ctx, username, password, companyID, apiURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ Helisa connection test failed")
		return err
	}

	uc.log.Info(ctx).Msg("✅ Helisa connection test successful")
	return nil
}
