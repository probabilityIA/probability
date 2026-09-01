package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	invoicing := router.Group("/invoicing")
	invoicing.Use(middleware.JWT())
	if h.moduleAccessMW != nil {
		invoicing.Use(h.moduleAccessMW)
	}
	{
		invoices := invoicing.Group("/invoices")
		{
			invoices.POST("", middleware.JWT(), h.CreateInvoice)
			invoices.POST("/manual", middleware.JWT(), h.RegisterManualInvoice)
			invoices.GET("", middleware.JWT(), h.ListInvoices)
			invoices.GET("/:id", middleware.JWT(), h.GetInvoice)
			invoices.POST("/:id/cancel", middleware.JWT(), h.CancelInvoice)
			invoices.POST("/:id/mark-as-cancelled", middleware.JWT(), h.MarkInvoiceAsCancelled)
			invoices.POST("/:id/refresh", middleware.JWT(), h.RefreshInvoice)
			invoices.POST("/:id/retry", middleware.JWT(), h.RetryInvoice)
			invoices.DELETE("/:id/retry", middleware.JWT(), h.CancelRetry)
			invoices.PUT("/:id/retry", middleware.JWT(), h.EnableRetry)
			invoices.POST("/:id/cash-receipt", middleware.JWT(), h.GenerateCashReceipt)
			invoices.DELETE("/:id", middleware.JWT(), h.DeleteInvoice)
			invoices.GET("/:id/sync-logs", middleware.JWT(), h.GetInvoiceSyncLogs)
			invoices.POST("/:id/credit-notes", middleware.JWT(), h.CreateCreditNote)

			invoices.GET("/invoiceable-orders", middleware.JWT(), h.ListInvoiceableOrders)
			invoices.POST("/bulk", middleware.JWT(), h.BulkCreateInvoices)
			invoices.POST("/retry-failed", middleware.JWT(), h.RetryFailedInvoices)

			invoices.POST("/compare", middleware.JWT(), h.CompareInvoices)
			invoices.GET("/compare/:correlationId", middleware.JWT(), h.GetCompareResult)
			invoices.POST("/sync-cancellations", middleware.JWT(), h.SyncCancellations)
			invoices.POST("/items", middleware.JWT(), h.ListItems)
			invoices.GET("/items/:correlationId", middleware.JWT(), h.GetListItemsResult)
			invoices.POST("/bank-accounts", middleware.JWT(), h.ListBankAccounts)
			invoices.GET("/bank-accounts/:correlationId", middleware.JWT(), h.GetListBankAccountsResult)
			invoices.GET("/cod-report", middleware.JWT(), h.GetCODReport)
			invoices.GET("/cod-report/invoices", middleware.JWT(), h.GetCODReportInvoices)
		}

		invoicing.POST("/journals", middleware.JWT(), h.CreateJournal)

		invoicing.POST("/inventory/sync", middleware.JWT(), h.SyncInventory)
		invoicing.POST("/inventory/siigo-warehouses", middleware.JWT(), h.ListSiigoWarehouses)

		providers := invoicing.Group("/providers")
		{
			providers.POST("", middleware.JWT(), h.CreateProvider)
			providers.GET("", middleware.JWT(), h.ListProviders)
			providers.GET("/:id", middleware.JWT(), h.GetProvider)
			providers.PUT("/:id", middleware.JWT(), h.UpdateProvider)
			providers.POST("/:id/test", middleware.JWT(), h.TestProvider)
		}

		configs := invoicing.Group("/configs")
		{
			configs.POST("", middleware.JWT(), h.CreateConfig)
			configs.GET("", middleware.JWT(), h.ListConfigs)
			configs.GET("/:id", middleware.JWT(), h.GetConfig)
			configs.PUT("/:id", middleware.JWT(), h.UpdateConfig)
			configs.DELETE("/:id", middleware.JWT(), h.DeleteConfig)
			configs.POST("/:id/enable", middleware.JWT(), h.EnableConfig)
			configs.POST("/:id/disable", middleware.JWT(), h.DisableConfig)
			configs.POST("/:id/enable-auto-invoice", middleware.JWT(), h.EnableAutoInvoice)
			configs.POST("/:id/disable-auto-invoice", middleware.JWT(), h.DisableAutoInvoice)
		}

		invoicing.GET("/summary", middleware.JWT(), h.GetSummary)
		invoicing.GET("/stats", middleware.JWT(), h.GetStats)
		invoicing.GET("/trends", middleware.JWT(), h.GetTrends)

		bulkJobs := invoicing.Group("/bulk-jobs")
		{
			bulkJobs.GET("", middleware.JWT(), h.ListBulkJobs)
			bulkJobs.GET("/:id", middleware.JWT(), h.GetBulkJobStatus)
		}
	}
}
