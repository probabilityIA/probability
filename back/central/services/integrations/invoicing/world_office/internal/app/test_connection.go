package app

import (
	"context"
	"fmt"
)

// TestConnection valida que las credenciales y configuración provistos sean correctos
// contra la API de World Office. Implementa ports.IInvoiceUseCase.TestConnection.
func (uc *invoicingUseCase) TestConnection(ctx context.Context, _ map[string]interface{}, credentials map[string]interface{}) error {
	uc.log.Info(ctx).Msg("🧪 Testing connection with World Office API")

	username, okUsername := credentials["username"].(string)
	password, okPassword := credentials["password"].(string)
	companyCode, okCompanyCode := credentials["company_code"].(string)
	apiURL, _ := credentials["api_url"].(string) // opcional

	uc.log.Info(ctx).
		Bool("has_username", okUsername && username != "").
		Bool("has_password", okPassword && password != "").
		Bool("has_company_code", okCompanyCode && companyCode != "").
		Bool("has_api_url", apiURL != "").
		Msg("📋 World Office credentials validation")

	if !okUsername || username == "" {
		return fmt.Errorf("el campo username es requerido en las credenciales")
	}
	if !okPassword || password == "" {
		return fmt.Errorf("el campo password es requerido en las credenciales")
	}
	if !okCompanyCode || companyCode == "" {
		return fmt.Errorf("el campo company_code es requerido en las credenciales")
	}

	uc.log.Info(ctx).Msg("🔌 Calling World Office client.TestAuthentication...")
	if err := uc.worldOfficeClient.TestAuthentication(ctx, username, password, companyCode, apiURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("❌ World Office connection test failed")
		return err
	}

	uc.log.Info(ctx).Msg("✅ World Office connection test successful")
	return nil
}
