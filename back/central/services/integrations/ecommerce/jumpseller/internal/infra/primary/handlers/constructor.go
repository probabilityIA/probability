package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	HandleWebhook(c *gin.Context)
	SyncProducts(c *gin.Context)
	ReconcileProducts(c *gin.Context)
	ApplyProducts(c *gin.Context)
	AssociateProducts(c *gin.Context)
	SyncInventory(c *gin.Context)
	GetLocations(c *gin.Context)
	InitiateOAuth(c *gin.Context)
	OAuthCallback(c *gin.Context)
	GetOAuthToken(c *gin.Context)
	VerifyApp(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type jumpsellerHandler struct {
	useCase         usecases.IJumpsellerUseCase
	coreIntegration integrationcore.IIntegrationCore
	config          env.IConfig
	logger          log.ILogger
}

func New(useCase usecases.IJumpsellerUseCase, coreIntegration integrationcore.IIntegrationCore, config env.IConfig, logger log.ILogger) IHandler {
	return &jumpsellerHandler{
		useCase:         useCase,
		coreIntegration: coreIntegration,
		config:          config,
		logger:          logger.WithModule("jumpseller"),
	}
}

func (h *jumpsellerHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	jumpseller := router.Group("/jumpseller")
	{
		jumpseller.POST("/webhook", h.HandleWebhook)
		jumpseller.GET("/callback", h.OAuthCallback)
		jumpseller.POST("/products/sync", middleware.JWT(), h.SyncProducts)
		jumpseller.POST("/products/reconcile", middleware.JWT(), h.ReconcileProducts)
		jumpseller.POST("/products/apply", middleware.JWT(), h.ApplyProducts)
		jumpseller.POST("/products/associate", middleware.JWT(), h.AssociateProducts)
		jumpseller.POST("/inventory/sync", middleware.JWT(), h.SyncInventory)
		jumpseller.GET("/locations", middleware.JWT(), h.GetLocations)
	}

	oauth := router.Group("/integrations/jumpseller")
	{
		oauth.POST("/connect", middleware.JWT(), h.InitiateOAuth)
		oauth.GET("/verify-app", middleware.JWT(), h.VerifyApp)
		oauth.GET("/oauth/token", h.GetOAuthToken)
	}
}
