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

		// Proveedores de facturación
		providers := invoicing.Group("/providers")
		{
			providers.POST("", h.CreateProvider)           // Crear proveedor
			providers.GET("", h.ListProviders)             // Listar proveedores
			providers.GET("/:id", h.GetProvider)           // Obtener proveedor
			providers.PUT("/:id", h.UpdateProvider)        // Actualizar proveedor
			providers.POST("/:id/test", h.TestProvider)    // Probar conexión
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
	}
}
