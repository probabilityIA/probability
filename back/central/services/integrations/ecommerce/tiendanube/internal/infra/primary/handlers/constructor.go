package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	HandleWebhook(c *gin.Context)
	SyncProducts(c *gin.Context)
	ReconcileProducts(c *gin.Context)
	ApplyProducts(c *gin.Context)
	AssociateProducts(c *gin.Context)
	SyncInventory(c *gin.Context)
	CompareInventory(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type tiendanubeHandler struct {
	useCase usecases.ITiendanubeUseCase
	logger  log.ILogger
}

func New(useCase usecases.ITiendanubeUseCase, logger log.ILogger) IHandler {
	return &tiendanubeHandler{
		useCase: useCase,
		logger:  logger.WithModule("tiendanube"),
	}
}

func (h *tiendanubeHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	tn := router.Group("/tiendanube")
	{
		tn.POST("/webhook", h.HandleWebhook)
		tn.POST("/products/sync", middleware.JWT(), h.SyncProducts)
		tn.POST("/products/reconcile", middleware.JWT(), h.ReconcileProducts)
		tn.POST("/products/apply", middleware.JWT(), h.ApplyProducts)
		tn.POST("/products/associate", middleware.JWT(), h.AssociateProducts)
		tn.POST("/inventory/sync", middleware.JWT(), h.SyncInventory)
		tn.POST("/inventory/compare", middleware.JWT(), h.CompareInventory)
	}
}
