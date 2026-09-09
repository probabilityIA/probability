package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	sf := router.Group("/storefront")
	{
		sf.POST("/register", h.Register)

		sf.GET("/catalog", middleware.JWT(), h.ListCatalog)
		sf.GET("/catalog/filters", middleware.JWT(), h.GetCatalogFilters)
		sf.PUT("/catalog/layout", middleware.JWT(), h.UpdateCatalogLayout)
		sf.PUT("/catalog/banner", middleware.JWT(), h.UpdateCatalogBanner)
		sf.POST("/catalog/banner/image", middleware.JWT(), h.UploadCatalogBannerImage)
		sf.GET("/catalog/:id", middleware.JWT(), h.GetProduct)
		sf.POST("/orders", middleware.JWT(), h.CreateOrder)
		sf.GET("/orders", middleware.JWT(), h.ListMyOrders)
		sf.GET("/orders/:id", middleware.JWT(), h.GetMyOrder)

		sf.POST("/clients", middleware.JWT(), h.CreateClient)
		sf.GET("/clients", middleware.JWT(), h.ListClients)
	}
}
