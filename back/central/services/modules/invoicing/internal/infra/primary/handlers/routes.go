package handlers

import "github.com/gin-gonic/gin"

// RegisterRoutes registra todas las rutas del módulo de facturación
func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	invoicing := router.Group("/invoicing")
	{
		// Facturas
		invoices := invoicing.Group("/invoices")
		{
			invoices.POST("", h.CreateInvoice)           // Crear factura manual
			invoices.GET("", h.ListInvoices)             // Listar facturas
			invoices.GET("/:id", h.GetInvoice)           // Obtener factura
			invoices.POST("/:id/cancel", h.CancelInvoice) // Cancelar factura
			invoices.POST("/:id/retry", h.RetryInvoice)   // Reintentar factura
			invoices.POST("/:id/credit-notes", h.CreateCreditNote) // Crear nota de crédito
		}

		// Proveedores de facturación (DEPRECADO - Migrado a integrations/core)
		// NOTA: Estas rutas están deprecadas y serán eliminadas en una futura versión
		// Usar endpoints de integrations/core para gestión de proveedores de facturación
		providers := invoicing.Group("/providers")
		{
			providers.POST("", h.CreateProvider)           // DEPRECATED: Crear proveedor
			providers.GET("", h.ListProviders)             // DEPRECATED: Listar proveedores
			providers.GET("/:id", h.GetProvider)           // DEPRECATED: Obtener proveedor
			providers.PUT("/:id", h.UpdateProvider)        // DEPRECATED: Actualizar proveedor
			providers.POST("/:id/test", h.TestProvider)    // DEPRECATED: Probar conexión
		}

		// Configuraciones de facturación
		configs := invoicing.Group("/configs")
		{
			configs.POST("", h.CreateConfig)       // Crear configuración
			configs.GET("", h.ListConfigs)         // Listar configuraciones
			configs.GET("/:id", h.GetConfig)       // Obtener configuración
			configs.PUT("/:id", h.UpdateConfig)    // Actualizar configuración
			configs.DELETE("/:id", h.DeleteConfig) // Eliminar configuración
		}

		// Estadísticas y resúmenes (NUEVO)
		invoicing.GET("/summary", h.GetSummary) // Resumen general con KPIs
		invoicing.GET("/stats", h.GetStats)     // Estadísticas detalladas
		invoicing.GET("/trends", h.GetTrends)   // Tendencias temporales
	}
}
