package businesshandler

import (
	"github.com/secamc93/probability/back/central/services/auth/middleware"

	"github.com/gin-gonic/gin"
)

func (h *BusinessHandler) RegisterRoutes(router *gin.RouterGroup, handler IBusinessHandler) {
	businesses := router.Group("/businesses")

	businesses.GET("", middleware.JWT(), handler.GetBusinesses)
	businesses.GET("/simple", middleware.JWT(), handler.GetBusinessesSimple)
	businesses.GET("/configured-resources", middleware.JWT(), handler.GetBusinessesConfiguredResourcesHandler)
	businesses.GET("/:id/configured-resources", middleware.JWT(), handler.GetBusinessConfiguredResourcesByIDHandler)
	businesses.GET("/:id", middleware.JWT(), handler.GetBusinessByIDHandler)
	businesses.POST("", middleware.JWT(), middleware.RequireSuperAdmin(), handler.CreateBusinessHandler)
	businesses.PUT("/:id", middleware.JWT(), handler.UpdateBusinessHandler)
	businesses.DELETE("/:id", middleware.JWT(), middleware.RequireSuperAdmin(), handler.DeleteBusinessHandler)

	businesses.PUT("/configured-resources/:resource_id/activate", middleware.JWT(), middleware.RequireSuperAdmin(), handler.ActivateBusinessResourceHandler)
	businesses.PUT("/configured-resources/:resource_id/deactivate", middleware.JWT(), middleware.RequireSuperAdmin(), handler.DeactivateBusinessResourceHandler)

	businesses.PUT("/:id/activate", middleware.JWT(), middleware.RequireSuperAdmin(), handler.ActivateBusinessHandler)
	businesses.PUT("/:id/deactivate", middleware.JWT(), middleware.RequireSuperAdmin(), handler.DeactivateBusinessHandler)
}
