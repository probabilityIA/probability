package businesshandler

import (
	"github.com/secamc93/probability/back/central/services/auth/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas del handler de Business
func (h *BusinessHandler) RegisterRoutes(router *gin.RouterGroup, handler IBusinessHandler) {
	businesses := router.Group("/businesses")

	// Rutas de Business
	businesses.GET("", middleware.JWT(), handler.GetBusinesses)
	businesses.GET("/simple", middleware.JWT(), handler.GetBusinessesSimple)
	businesses.GET("/configured-resources", middleware.JWT(), handler.GetBusinessesConfiguredResourcesHandler)
	businesses.GET("/:id/configured-resources", middleware.JWT(), handler.GetBusinessConfiguredResourcesByIDHandler)
	businesses.GET("/:id", middleware.JWT(), handler.GetBusinessByIDHandler)
	businesses.POST("", middleware.JWT(), handler.CreateBusinessHandler)
	businesses.PUT("/:id", middleware.JWT(), handler.UpdateBusinessHandler)
	businesses.DELETE("/:id", middleware.JWT(), handler.DeleteBusinessHandler)

	// Rutas para activar/desactivar recursos de business
	businesses.PUT("/configured-resources/:resource_id/activate", middleware.JWT(), handler.ActivateBusinessResourceHandler)
	businesses.PUT("/configured-resources/:resource_id/deactivate", middleware.JWT(), handler.DeactivateBusinessResourceHandler)

	// Rutas para activar/desactivar business
	businesses.PUT("/:id/activate", middleware.JWT(), handler.ActivateBusinessHandler)
	businesses.PUT("/:id/deactivate", middleware.JWT(), handler.DeactivateBusinessHandler)
}
