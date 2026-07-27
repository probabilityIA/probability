package usecases

import (
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func normalizeBaseURL(raw string) string {
	baseURL := strings.TrimSpace(raw)
	if baseURL == "" {
		return ""
	}
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	return strings.TrimRight(baseURL, "/")
}

func resolveEffectiveBaseURL(integration *domain.Integration) string {
	if integration == nil || !integration.IsTesting {
		return ""
	}
	return normalizeBaseURL(integration.BaseURLTest)
}
