package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	core "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	HandleNotification(c *gin.Context)
	InitiateOAuthHandler(c *gin.Context)
	OAuthCallbackHandler(c *gin.Context)
	GetOAuthTokenHandler(c *gin.Context)
	ReconcileProducts(c *gin.Context)
	ApplyProducts(c *gin.Context)
	SyncInventory(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type meliHandler struct {
	useCase         usecases.IMeliUseCase
	logger          log.ILogger
	config          env.IConfig
	coreIntegration core.IIntegrationCore
}

func New(useCase usecases.IMeliUseCase, logger log.ILogger, config env.IConfig, coreIntegration core.IIntegrationCore) IHandler {
	return &meliHandler{
		useCase:         useCase,
		logger:          logger.WithModule("meli"),
		config:          config,
		coreIntegration: coreIntegration,
	}
}

func (h *meliHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	meli := router.Group("/meli")
	{
		meli.POST("/notifications", h.HandleNotification)
	}

	oauthGroup := router.Group("/integrations/meli")
	{
		oauthGroup.POST("/connect", middleware.JWT(), h.InitiateOAuthHandler)
		oauthGroup.GET("/oauth/token", h.GetOAuthTokenHandler)
		oauthGroup.POST("/products/reconcile", middleware.JWT(), h.ReconcileProducts)
		oauthGroup.POST("/products/apply", middleware.JWT(), h.ApplyProducts)
		oauthGroup.POST("/products/associate", middleware.JWT(), h.AssociateProducts)
		oauthGroup.POST("/inventory/sync", middleware.JWT(), h.SyncInventory)
	}

	router.GET("/meli/callback", h.OAuthCallbackHandler)
}
