package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiktok/internal/app/usecases"
)

type TikTokCore struct {
	integrationcore.BaseIntegration
	useCase usecases.ITikTokUseCase
}

func New(useCase usecases.ITikTokUseCase) *TikTokCore {
	return &TikTokCore{useCase: useCase}
}

func (t *TikTokCore) TestConnection(ctx context.Context, config map[string]interface{}, credentials map[string]interface{}) error {
	return t.useCase.TestConnection(ctx, config, credentials)
}
