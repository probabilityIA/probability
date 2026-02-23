package handlers

import (
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ShopifyHandler struct {
	useCase usecases.IShopifyUseCase
	logger  log.ILogger
	config  env.IConfig
}

func New(useCase usecases.IShopifyUseCase, logger log.ILogger, config env.IConfig) *ShopifyHandler {
	contextualLogger := logger.WithModule("shopify")
	return &ShopifyHandler{
		useCase: useCase,
		logger:  contextualLogger,
		config:  config,
	}
}
