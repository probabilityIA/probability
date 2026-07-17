package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/log"
)

type IHandler interface {
	HandleWebhook(c *gin.Context)
	GetWebhookStatus(c *gin.Context)
	RegisterWebhook(c *gin.Context)
	UnregisterWebhook(c *gin.Context)
	SyncProducts(c *gin.Context)
	ReconcileProducts(c *gin.Context)
	ApplyProducts(c *gin.Context)
	AssociateProducts(c *gin.Context)
	SyncInventory(c *gin.Context)
	GetWarehouses(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup, logger log.ILogger)
}

type vtexHandler struct {
	useCase usecases.IVTEXUseCase
	baseURL string
	logger  log.ILogger
}

func New(useCase usecases.IVTEXUseCase, baseURL string, logger log.ILogger) IHandler {
	return &vtexHandler{
		useCase: useCase,
		baseURL: baseURL,
		logger:  logger.WithModule("vtex"),
	}
}

func (h *vtexHandler) RegisterRoutes(router *gin.RouterGroup, logger log.ILogger) {
	vtex := router.Group("/vtex")
	{
		vtex.POST("/webhook", h.HandleWebhook)
		vtex.GET("/webhooks", middleware.JWT(), h.GetWebhookStatus)
		vtex.POST("/webhooks/register", middleware.JWT(), h.RegisterWebhook)
		vtex.POST("/webhooks/unregister", middleware.JWT(), h.UnregisterWebhook)
		vtex.POST("/products/sync", middleware.JWT(), h.SyncProducts)
		vtex.POST("/products/reconcile", middleware.JWT(), h.ReconcileProducts)
		vtex.POST("/products/apply", middleware.JWT(), h.ApplyProducts)
		vtex.POST("/products/associate", middleware.JWT(), h.AssociateProducts)
		vtex.POST("/inventory/sync", middleware.JWT(), h.SyncInventory)
		vtex.GET("/warehouses", middleware.JWT(), h.GetWarehouses)
	}
}
