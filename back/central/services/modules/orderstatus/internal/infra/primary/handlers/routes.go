package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo de order status mappings
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	// Rutas para estados de órdenes
	statuses := router.Group("/order-statuses")
	{
		statuses.GET("", h.ListOrderStatuses)              // GET /api/v1/order-statuses
		statuses.GET("/simple", h.ListOrderStatusesSimple) // GET /api/v1/order-statuses/simple
		statuses.POST("", h.CreateOrderStatus)             // POST /api/v1/order-statuses
		statuses.GET("/:id", h.GetOrderStatus)             // GET /api/v1/order-statuses/:id
		statuses.PUT("/:id", h.UpdateOrderStatus)          // PUT /api/v1/order-statuses/:id
		statuses.DELETE("/:id", h.DeleteOrderStatus)       // DELETE /api/v1/order-statuses/:id
	}

	// Rutas para tipos de integración ecommerce (scope-aware)
	router.GET("/ecommerce-integration-types", h.ListEcommerceIntegrationTypes) // GET /api/v1/ecommerce-integration-types

	// Rutas para estados de canales de integración
	channelStatuses := router.Group("/channel-statuses")
	{
		channelStatuses.GET("", h.ListChannelStatuses)      // GET /api/v1/channel-statuses?integration_type_id=1
		channelStatuses.POST("", h.CreateChannelStatus)     // POST /api/v1/channel-statuses
		channelStatuses.PUT("/:id", h.UpdateChannelStatus)  // PUT /api/v1/channel-statuses/:id
		channelStatuses.DELETE("/:id", h.DeleteChannelStatus) // DELETE /api/v1/channel-statuses/:id
	}

	// Rutas para estados de fulfillment
	fulfillmentStatuses := router.Group("/fulfillment-statuses")
	{
		fulfillmentStatuses.GET("", h.ListFulfillmentStatuses) // GET /api/v1/fulfillment-statuses
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
