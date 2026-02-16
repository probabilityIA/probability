package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra todas las rutas del módulo shipments
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	shipments := router.Group("/shipments")
	{
		// CRUD básico
		shipments.GET("", h.ListShipments)
		shipments.GET("/:id", h.GetShipmentByID)
		shipments.POST("", h.CreateShipment)
		shipments.PUT("/:id", h.UpdateShipment)
		shipments.DELETE("/:id", h.DeleteShipment)

		// Rutas adicionales
		shipments.GET("/order/:order_id", h.GetShipmentsByOrderID)
		shipments.GET("/tracking/:tracking_number", h.GetShipmentByTrackingNumber)

		// Direcciones de Origen
		shipments.GET("/origin-addresses", h.ListOriginAddresses)
		shipments.POST("/origin-addresses", h.CreateOriginAddress)
		shipments.PUT("/origin-addresses/:id", h.UpdateOriginAddress)
		shipments.DELETE("/origin-addresses/:id", h.DeleteOriginAddress)

		// EnvioClick Integration
		shipments.POST("/quote", h.QuoteShipment)
		shipments.POST("/generate", h.GenerateGuide)
		shipments.POST("/tracking/:tracking_number/track", h.TrackEnvioclickShipment)
		shipments.POST("/:id/cancel", h.CancelEnvioclickShipment)
	}

}
