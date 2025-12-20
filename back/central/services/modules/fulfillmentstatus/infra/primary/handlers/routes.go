package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo de fulfillment statuses
func (h *FulfillmentStatusHandlers) RegisterRoutes(router *gin.RouterGroup) {
	// Rutas para estados de fulfillment
	statuses := router.Group("/fulfillment-statuses")
	{
		statuses.GET("", h.ListFulfillmentStatuses) // GET /api/v1/fulfillment-statuses
	}
}
