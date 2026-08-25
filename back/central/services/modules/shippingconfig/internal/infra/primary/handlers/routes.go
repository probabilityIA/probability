package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/shipping-config")
	group.Use(middleware.JWT())
	{
		group.GET("", h.GetOverview)
		group.PUT("", h.SaveBusinessConfig)
		group.PUT("/warehouses/:warehouse_id", h.SaveWarehouseConfig)
		group.DELETE("/warehouses/:warehouse_id", h.DeleteWarehouseConfig)
		group.PUT("/warehouses/:warehouse_id/default", h.SetDefaultWarehouse)
	}
}
