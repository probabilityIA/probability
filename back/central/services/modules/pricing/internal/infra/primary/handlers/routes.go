package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	pricing := router.Group("/pricing")
	{
		groups := pricing.Group("/client-groups")
		groups.GET("", middleware.JWT(), h.ListClientGroups)
		groups.POST("", middleware.JWT(), h.CreateClientGroup)
		groups.GET("/:id", middleware.JWT(), h.GetClientGroup)
		groups.PUT("/:id", middleware.JWT(), h.UpdateClientGroup)
		groups.DELETE("/:id", middleware.JWT(), h.DeleteClientGroup)
		groups.GET("/:id/members", middleware.JWT(), h.ListGroupMembers)
		groups.POST("/:id/members", middleware.JWT(), h.AddGroupMembers)
		groups.DELETE("/:id/members/:clientId", middleware.JWT(), h.RemoveGroupMember)

		pricing.GET("/clients", middleware.JWT(), h.ListAvailableClients)
		pricing.GET("/catalog-prices", middleware.JWT(), h.GetCatalogPrices)
		pricing.PUT("/catalog-prices", middleware.JWT(), h.SaveCatalogPrices)
		pricing.GET("/effective-price", middleware.JWT(), h.GetEffectivePrice)
	}
}
