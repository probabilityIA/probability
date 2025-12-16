package handlerintegrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IIntegrationHandler interface {
	GetIntegrationsHandler(c *gin.Context)
	GetIntegrationByIDHandler(c *gin.Context)
	GetIntegrationByTypeHandler(c *gin.Context)
	CreateIntegrationHandler(c *gin.Context)
	UpdateIntegrationHandler(c *gin.Context)
	DeleteIntegrationHandler(c *gin.Context)
	TestIntegrationHandler(c *gin.Context)
	TestConnectionRawHandler(c *gin.Context)
	SyncOrdersByIntegrationIDHandler(c *gin.Context)
	SyncOrdersByBusinessHandler(c *gin.Context)
	ActivateIntegrationHandler(c *gin.Context)
	DeactivateIntegrationHandler(c *gin.Context)
	SetAsDefaultHandler(c *gin.Context)
	GetWebhookURLHandler(c *gin.Context)
	ListWebhooksHandler(c *gin.Context)
	DeleteWebhookHandler(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}
type IntegrationHandler struct {
	usecase      usecaseintegrations.IIntegrationUseCase
	logger       log.ILogger
	orderSyncSvc domain.IOrderSyncService
	env          env.IConfig
}

func New(usecase usecaseintegrations.IIntegrationUseCase, logger log.ILogger, orderSyncSvc domain.IOrderSyncService, env env.IConfig) IIntegrationHandler {
	contextualLogger := logger.WithModule("integrations")
	return &IntegrationHandler{
		usecase:      usecase,
		logger:       contextualLogger,
		orderSyncSvc: orderSyncSvc,
		env:          env,
	}
}

// getImageURLBase obtiene la URL base de S3 para construir URLs completas
func (h *IntegrationHandler) getImageURLBase() string {
	return h.env.Get("URL_BASE_DOMAIN_S3")
}
