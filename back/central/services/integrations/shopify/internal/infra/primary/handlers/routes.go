package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *ShopifyHandlers) RegisterRoutes(router *gin.RouterGroup) {
	shopifyGroup := router.Group("/shopify")
	{
		shopifyGroup.POST("/sync/:id", h.SyncOrders)
	}
}
