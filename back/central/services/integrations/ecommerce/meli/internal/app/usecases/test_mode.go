package usecases

import (
	"context"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
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

func (uc *meliUseCase) clientFor(ctx context.Context, integration *domain.Integration) domain.IMeliClient {
	cli := uc.client.ForAccount(strconv.FormatUint(uint64(integration.ID), 10))

	baseURL := resolveEffectiveBaseURL(integration)
	if baseURL == "" {
		return cli
	}
	uc.logger.Info(ctx).
		Uint("integration_id", integration.ID).
		Str("base_url", baseURL).
		Msg("Meli en modo pruebas: usando base_url_test")
	return cli.WithBaseURL(baseURL)
}
