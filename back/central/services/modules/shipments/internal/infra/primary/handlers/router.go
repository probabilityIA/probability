package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo shipments
func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/tracking/search", h.PublicSearchTracking)
	router.GET("/tracking/:tracking_number/history", h.PublicGetTrackingHistory)

	router.POST("/shopify/shipping-rates/:integration_id", h.ShopifyShippingRates)

	shipments := router.Group("/shipments", middleware.JWT())
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
		shipments.GET("/quotes", h.ListSavedQuotes)
		shipments.GET("/quotes/:id", h.GetSavedQuote)

		shipments.GET("/origin-addresses", h.ListOriginAddresses)
		shipments.POST("/origin-addresses", h.CreateOriginAddress)
		shipments.PUT("/origin-addresses/:id", h.UpdateOriginAddress)
		shipments.DELETE("/origin-addresses/:id", h.DeleteOriginAddress)

		// Transport Operations (carrier-agnostic — resolved dynamically per business)
		shipments.POST("/quote", h.QuoteShipment)
		shipments.POST("/generate", h.GenerateGuide)
		shipments.GET("/guide-formats", h.ListGuideFormats)
		shipments.GET("/:id/guide", h.RenderGuide)
		shipments.POST("/:id/extract-coordinadora-data", h.ExtractCoordinadoraData)
		shipments.POST("/tracking/:tracking_number/track", h.TrackShipment)
		shipments.POST("/:id/cancel", h.CancelShipment)
		shipments.POST("/cancel-batch", h.CancelBatchShipments)
		shipments.POST("/sync-status", h.SyncShipmentStatus)

		shipments.GET("/stats/by-geozone", h.StatsByGeozone)

		shipments.GET("/cod", h.ListCODShipments)
		shipments.POST("/:id/collect-cod", h.CollectCOD)

		shipments.GET("/manifest/carriers", h.ListManifestCarriers)
		shipments.GET("/manifest/pending", h.ListManifestPending)
		shipments.POST("/manifest/pdf", h.GenerateManifestPDF)
	}

}
