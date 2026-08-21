package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	HandleWebhook(c *gin.Context)
	HandleStoreRedact(c *gin.Context)
	HandleCustomersRedact(c *gin.Context)
	HandleCustomersDataRequest(c *gin.Context)
	SyncProducts(c *gin.Context)
	ReconcileProducts(c *gin.Context)
	ApplyProducts(c *gin.Context)
	AssociateProducts(c *gin.Context)
	SyncInventory(c *gin.Context)
	CompareInventory(c *gin.Context)
	InitiateOAuth(c *gin.Context)
	OAuthCallback(c *gin.Context)
	GetOAuthToken(c *gin.Context)
	VerifyApp(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type tiendanubeHandler struct {
	useCase         usecases.ITiendanubeUseCase
	coreIntegration integrationcore.IIntegrationCore
	config          env.IConfig
	logger          log.ILogger
}

func New(useCase usecases.ITiendanubeUseCase, coreIntegration integrationcore.IIntegrationCore, config env.IConfig, logger log.ILogger) IHandler {
	return &tiendanubeHandler{
		useCase:         useCase,
		coreIntegration: coreIntegration,
		config:          config,
		logger:          logger.WithModule("tiendanube"),
	}
}

func (h *tiendanubeHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	tn := router.Group("/tiendanube")
	{
		tn.POST("/webhook", h.HandleWebhook)
		tn.GET("/callback", h.OAuthCallback)
		tn.POST("/webhook/store-redact", h.HandleStoreRedact)
		tn.POST("/webhook/customers-redact", h.HandleCustomersRedact)
		tn.POST("/webhook/customers-data-request", h.HandleCustomersDataRequest)
		tn.POST("/products/sync", middleware.JWT(), h.SyncProducts)
		tn.POST("/products/reconcile", middleware.JWT(), h.ReconcileProducts)
		tn.POST("/products/apply", middleware.JWT(), h.ApplyProducts)
		tn.POST("/products/associate", middleware.JWT(), h.AssociateProducts)
		tn.POST("/inventory/sync", middleware.JWT(), h.SyncInventory)
		tn.POST("/inventory/compare", middleware.JWT(), h.CompareInventory)
	}

	oauth := router.Group("/integrations/tiendanube")
	{
		oauth.POST("/connect", middleware.JWT(), h.InitiateOAuth)
		oauth.GET("/verify-app", middleware.JWT(), h.VerifyApp)
		oauth.GET("/oauth/token", h.GetOAuthToken)
	}
}
