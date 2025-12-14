package handlerintegrations

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/app/usecaseintegrations"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
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
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}
type IntegrationHandler struct {
	usecase      usecaseintegrations.IIntegrationUseCase
	logger       log.ILogger
	orderSyncSvc domain.IOrderSyncService
}

func New(usecase usecaseintegrations.IIntegrationUseCase, logger log.ILogger, orderSyncSvc domain.IOrderSyncService) IIntegrationHandler {
	contextualLogger := logger.WithModule("integrations")
	return &IntegrationHandler{
		usecase:      usecase,
		logger:       contextualLogger,
		orderSyncSvc: orderSyncSvc,
	}
}
