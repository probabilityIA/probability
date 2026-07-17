package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/app/usecases"
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
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type jumpsellerHandler struct {
	useCase usecases.IJumpsellerUseCase
	logger  log.ILogger
}

func New(useCase usecases.IJumpsellerUseCase, logger log.ILogger) IHandler {
	return &jumpsellerHandler{
		useCase: useCase,
		logger:  logger.WithModule("jumpseller"),
	}
}

func (h *jumpsellerHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	jumpseller := router.Group("/jumpseller")
	{
		jumpseller.POST("/webhook", h.HandleWebhook)
		jumpseller.POST("/products/sync", middleware.JWT(), h.SyncProducts)
		jumpseller.POST("/products/reconcile", middleware.JWT(), h.ReconcileProducts)
		jumpseller.POST("/products/apply", middleware.JWT(), h.ApplyProducts)
		jumpseller.POST("/products/associate", middleware.JWT(), h.AssociateProducts)
		jumpseller.POST("/inventory/sync", middleware.JWT(), h.SyncInventory)
		jumpseller.GET("/locations", middleware.JWT(), h.GetLocations)
	}
}
