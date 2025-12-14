package handlers

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ShopifyHandler struct {
	useCase         usecases.IShopifyUseCase
	logger          log.ILogger
	coreIntegration core.IIntegrationCore
}

func New(useCase usecases.IShopifyUseCase, logger log.ILogger, coreIntegration core.IIntegrationCore) *ShopifyHandler {
	contextualLogger := logger.WithModule("shopify")
	return &ShopifyHandler{
		useCase:         useCase,
		logger:          contextualLogger,
		coreIntegration: coreIntegration,
	}
}
