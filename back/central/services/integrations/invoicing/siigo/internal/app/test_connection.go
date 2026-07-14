package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/entities"
)

func (uc *invoicingUseCase) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	uc.log.Info(ctx).Msg("Testing connection with Siigo API")

	username, _ := credentials["username"].(string)
	accessKey, _ := credentials["access_key"].(string)
	accountID, _ := credentials["account_id"].(string)
	partnerID, _ := credentials["partner_id"].(string)

	if username == "" {
		return fmt.Errorf("el campo username (email) es requerido en las credenciales")
	}
	if accessKey == "" {
		return fmt.Errorf("el campo access_key es requerido en las credenciales")
	}
	if partnerID == "" {
		return fmt.Errorf("el campo partner_id es requerido en las credenciales")
	}

	apiURL := resolveBaseURL(config, credentials)

	uc.log.Info(ctx).
		Bool("has_username", username != "").
		Bool("has_access_key", accessKey != "").
		Bool("has_account_id", accountID != "").
		Bool("has_partner_id", partnerID != "").
		Str("resolved_url", apiURL).
		Msg("Siigo credentials validation")

	if apiURL == "" {
		return fmt.Errorf("URL de Siigo no configurada en el tipo de integracion (base_url o base_url_test)")
	}

	if err := uc.siigoClient.TestAuthentication(ctx, username, accessKey, accountID, partnerID, apiURL); err != nil {
		uc.log.Error(ctx).Err(err).Msg("Siigo connection test failed")
		return err
	}

	uc.log.Info(ctx).Msg("Siigo connection test successful")
	return nil
}

func resolveBaseURL(config map[string]interface{}, credentials map[string]interface{}) string {
	var isTesting bool
	var baseURLTest, baseURL string
	if config != nil {
		isTesting, _ = config["is_testing"].(bool)
		baseURLTest, _ = config["base_url_test"].(string)
		baseURL, _ = config["base_url"].(string)
	}
	var apiURL string
	if credentials != nil {
		apiURL, _ = credentials["api_url"].(string)
	}
	return entities.ResolveSiigoBaseURL(isTesting, baseURLTest, apiURL, baseURL)
}
