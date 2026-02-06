package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

// RegisterRoutes registra todas las rutas del módulo de facturación
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	invoicing := router.Group("/invoicing")
	{
		// Facturas
		invoices := invoicing.Group("/invoices")
		{
			invoices.POST("", middleware.JWT(), h.CreateInvoice)           // Crear factura manual
			invoices.GET("", middleware.JWT(), h.ListInvoices)             // Listar facturas
			invoices.GET("/:id", middleware.JWT(), h.GetInvoice)           // Obtener factura
			invoices.POST("/:id/cancel", middleware.JWT(), h.CancelInvoice) // Cancelar factura
			invoices.POST("/:id/retry", middleware.JWT(), h.RetryInvoice)   // Reintentar factura
			invoices.POST("/:id/credit-notes", middleware.JWT(), h.CreateCreditNote) // Crear nota de crédito

			// Creación masiva de facturas
			invoices.GET("/invoiceable-orders", middleware.JWT(), h.ListInvoiceableOrders) // Listar órdenes facturables
			invoices.POST("/bulk", middleware.JWT(), h.BulkCreateInvoices)                 // Crear facturas masivamente
		}

		// Proveedores de facturación (DEPRECADO - Migrado a integrations/core)
		// NOTA: Estas rutas están deprecadas y serán eliminadas en una futura versión
		// Usar endpoints de integrations/core para gestión de proveedores de facturación
		providers := invoicing.Group("/providers")
		{
			providers.POST("", middleware.JWT(), h.CreateProvider)           // DEPRECATED: Crear proveedor
			providers.GET("", middleware.JWT(), h.ListProviders)             // DEPRECATED: Listar proveedores
			providers.GET("/:id", middleware.JWT(), h.GetProvider)           // DEPRECATED: Obtener proveedor
			providers.PUT("/:id", middleware.JWT(), h.UpdateProvider)        // DEPRECATED: Actualizar proveedor
			providers.POST("/:id/test", middleware.JWT(), h.TestProvider)    // DEPRECATED: Probar conexión
		}

		// Configuraciones de facturación
		configs := invoicing.Group("/configs")
		{
			configs.POST("", middleware.JWT(), h.CreateConfig)       // Crear configuración
			configs.GET("", middleware.JWT(), h.ListConfigs)         // Listar configuraciones
			configs.GET("/:id", middleware.JWT(), h.GetConfig)       // Obtener configuración
			configs.PUT("/:id", middleware.JWT(), h.UpdateConfig)    // Actualizar configuración
			configs.DELETE("/:id", middleware.JWT(), h.DeleteConfig) // Eliminar configuración
		}

		// Estadísticas y resúmenes (NUEVO)
		invoicing.GET("/summary", middleware.JWT(), h.GetSummary) // Resumen general con KPIs
		invoicing.GET("/stats", middleware.JWT(), h.GetStats)     // Estadísticas detalladas
		invoicing.GET("/trends", middleware.JWT(), h.GetTrends)   // Tendencias temporales
	}
}
