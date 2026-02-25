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
			invoices.POST("", middleware.JWT(), h.CreateInvoice)                     // Crear factura manual
			invoices.GET("", middleware.JWT(), h.ListInvoices)                       // Listar facturas
			invoices.GET("/:id", middleware.JWT(), h.GetInvoice)                     // Obtener factura
			invoices.POST("/:id/cancel", middleware.JWT(), h.CancelInvoice)          // Cancelar factura
			invoices.POST("/:id/retry", middleware.JWT(), h.RetryInvoice)            // Reintentar factura
			invoices.DELETE("/:id/retry", middleware.JWT(), h.CancelRetry)           // Cancelar reintentos pendientes
			invoices.PUT("/:id/retry", middleware.JWT(), h.EnableRetry)              // Habilitar reintentos automáticos
			invoices.GET("/:id/sync-logs", middleware.JWT(), h.GetInvoiceSyncLogs)   // Historial de sincronización
			invoices.POST("/:id/credit-notes", middleware.JWT(), h.CreateCreditNote) // Crear nota de crédito

			// Creación masiva de facturas
			invoices.GET("/invoiceable-orders", middleware.JWT(), h.ListInvoiceableOrders) // Listar órdenes facturables
			invoices.POST("/bulk", middleware.JWT(), h.BulkCreateInvoices)                 // Crear facturas masivamente

			// Comparación con proveedor (auditoría esporádica)
			invoices.POST("/compare", middleware.JWT(), h.CompareInvoices) // Iniciar comparación
		}

		// Proveedores de facturación (DEPRECADO - Migrado a integrations/core)
		// NOTA: Estas rutas están deprecadas y serán eliminadas en una futura versión
		// Usar endpoints de integrations/core para gestión de proveedores de facturación
		providers := invoicing.Group("/providers")
		{
			providers.POST("", middleware.JWT(), h.CreateProvider)        // DEPRECATED: Crear proveedor
			providers.GET("", middleware.JWT(), h.ListProviders)          // DEPRECATED: Listar proveedores
			providers.GET("/:id", middleware.JWT(), h.GetProvider)        // DEPRECATED: Obtener proveedor
			providers.PUT("/:id", middleware.JWT(), h.UpdateProvider)     // DEPRECATED: Actualizar proveedor
			providers.POST("/:id/test", middleware.JWT(), h.TestProvider) // DEPRECATED: Probar conexión
		}

		// Configuraciones de facturación
		configs := invoicing.Group("/configs")
		{
			configs.POST("", middleware.JWT(), h.CreateConfig)                                // Crear configuración
			configs.GET("", middleware.JWT(), h.ListConfigs)                                  // Listar configuraciones
			configs.GET("/:id", middleware.JWT(), h.GetConfig)                                // Obtener configuración
			configs.PUT("/:id", middleware.JWT(), h.UpdateConfig)                             // Actualizar configuración
			configs.DELETE("/:id", middleware.JWT(), h.DeleteConfig)                          // Eliminar configuración
			configs.POST("/:id/enable", middleware.JWT(), h.EnableConfig)                     // Activar configuración
			configs.POST("/:id/disable", middleware.JWT(), h.DisableConfig)                   // Desactivar configuración
			configs.POST("/:id/enable-auto-invoice", middleware.JWT(), h.EnableAutoInvoice)   // Activar facturación automática
			configs.POST("/:id/disable-auto-invoice", middleware.JWT(), h.DisableAutoInvoice) // Desactivar facturación automática
		}

		// Estadísticas y resúmenes (NUEVO)
		invoicing.GET("/summary", middleware.JWT(), h.GetSummary) // Resumen general con KPIs
		invoicing.GET("/stats", middleware.JWT(), h.GetStats)     // Estadísticas detalladas
		invoicing.GET("/trends", middleware.JWT(), h.GetTrends)   // Tendencias temporales

		// Jobs de facturación masiva (NUEVO - Asíncrono)
		bulkJobs := invoicing.Group("/bulk-jobs")
		{
			bulkJobs.GET("", middleware.JWT(), h.ListBulkJobs)         // Listar jobs
			bulkJobs.GET("/:id", middleware.JWT(), h.GetBulkJobStatus) // Estado de job
		}
	}
}
