package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo de order status mappings
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	// Rutas para estados de órdenes
	statuses := router.Group("/order-statuses")
	{
		statuses.GET("", h.ListOrderStatuses)              // GET /api/v1/order-statuses
		statuses.GET("/simple", h.ListOrderStatusesSimple) // GET /api/v1/order-statuses/simple
	}

	// Rutas para mapeos de estados
	mappings := router.Group("/order-status-mappings")
	{
		mappings.GET("", h.List)                // GET /api/v1/order-status-mappings
		mappings.GET("/:id", h.Get)             // GET /api/v1/order-status-mappings/:id
		mappings.POST("", h.Create)             // POST /api/v1/order-status-mappings
		mappings.PUT("/:id", h.Update)          // PUT /api/v1/order-status-mappings/:id
		mappings.DELETE("/:id", h.Delete)       // DELETE /api/v1/order-status-mappings/:id
		mappings.PATCH("/:id/toggle", h.Toggle) // PATCH /api/v1/order-status-mappings/:id/toggle
	}
}
