package handlers

import "github.com/gin-gonic/gin"

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	routes := router.Group("/routes")
	{
		routes.GET("", h.ListRoutes)
		routes.GET("/available-drivers", h.ListAvailableDrivers)
		routes.GET("/available-vehicles", h.ListAvailableVehicles)
		routes.GET("/assignable-orders", h.ListAssignableOrders)
		routes.GET("/:id", h.GetRoute)
		routes.POST("", h.CreateRoute)
		routes.PUT("/:id", h.UpdateRoute)
		routes.DELETE("/:id", h.DeleteRoute)

		// Lifecycle
		routes.POST("/:id/start", h.StartRoute)
		routes.POST("/:id/complete", h.CompleteRoute)

		// Stops
		routes.POST("/:id/stops", h.AddStop)
		routes.PUT("/:id/stops/:stopId", h.UpdateStop)
		routes.DELETE("/:id/stops/:stopId", h.DeleteStop)
		routes.POST("/:id/stops/:stopId/status", h.UpdateStopStatus)
		routes.PUT("/:id/stops/reorder", h.ReorderStops)
	}
}
