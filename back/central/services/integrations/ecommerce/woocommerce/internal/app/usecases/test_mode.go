package usecases

import "github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"

func resolveEffectiveStoreURL(integration *domain.Integration, storeURL string) string {
	if integration != nil && integration.IsTesting && integration.BaseURLTest != "" {
		return integration.BaseURLTest
	}
	return storeURL
}
