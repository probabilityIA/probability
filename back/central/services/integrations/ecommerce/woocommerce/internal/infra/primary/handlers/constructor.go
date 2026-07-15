package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	HandleWebhook(c *gin.Context)
	SyncProducts(c *gin.Context)
	ReconcileProducts(c *gin.Context)
	ApplyProducts(c *gin.Context)
	AssociateProducts(c *gin.Context)
	SyncInventory(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type wooCommerceHandler struct {
	useCase usecases.IWooCommerceUseCase
	logger  log.ILogger
}

func New(useCase usecases.IWooCommerceUseCase, logger log.ILogger) IHandler {
	return &wooCommerceHandler{
		useCase: useCase,
		logger:  logger.WithModule("woocommerce"),
	}
}

func (h *wooCommerceHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	woo := router.Group("/woocommerce")
	{
		woo.POST("/webhook", h.HandleWebhook)
		woo.POST("/products/sync", middleware.JWT(), h.SyncProducts)
		woo.POST("/products/reconcile", middleware.JWT(), h.ReconcileProducts)
		woo.POST("/products/apply", middleware.JWT(), h.ApplyProducts)
		woo.POST("/products/associate", middleware.JWT(), h.AssociateProducts)
		woo.POST("/inventory/sync", middleware.JWT(), h.SyncInventory)
	}
}
