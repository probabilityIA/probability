package handlers

import (
	core "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type ShopifyHandler struct {
	useCase         usecases.IShopifyUseCase
	logger          log.ILogger
	config          env.IConfig
	coreIntegration core.IIntegrationCore
	rabbit          rabbitmq.IQueue
}

func New(useCase usecases.IShopifyUseCase, logger log.ILogger, config env.IConfig, coreIntegration core.IIntegrationCore, rabbit rabbitmq.IQueue) *ShopifyHandler {
	contextualLogger := logger.WithModule("shopify")
	return &ShopifyHandler{
		useCase:         useCase,
		logger:          contextualLogger,
		config:          config,
		coreIntegration: coreIntegration,
		rabbit:          rabbit,
	}
}
