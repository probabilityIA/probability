package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
)

func (h *Handlers) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/accounting")
	g.Use(middleware.JWT())
	{
		g.GET("/concepts", h.ListConcepts)
		g.POST("/concepts", h.CreateConcept)
		g.PUT("/concepts/:id", h.UpdateConcept)
		g.PUT("/concepts/:id/taxes/:taxId", h.SetConceptTax)
		g.GET("/taxes", h.ListTaxes)
		g.POST("/taxes", h.CreateTax)
		g.PUT("/taxes/:id", h.UpdateTax)
		g.GET("/entries", h.ListEntries)
		g.POST("/entries", h.CreateEntry)
		g.DELETE("/entries/:id", h.DeleteEntry)
		g.GET("/report", h.Report)
		g.POST("/sync", h.Sync)
		g.GET("/invoices", h.ListInvoices)
		g.POST("/invoices", h.CreateInvoice)
		g.GET("/invoices/:id", h.GetInvoice)
		g.PUT("/invoices/:id", h.UpdateInvoice)
		g.DELETE("/invoices/:id", h.DeleteInvoice)
		g.POST("/invoices/:id/send", h.SendInvoice)
		g.POST("/invoices/:id/pay", h.MarkInvoicePaid)
		g.POST("/invoices/:id/cancel", h.CancelInvoice)
		g.POST("/invoices/:id/emit-dian", h.EmitInvoiceDian)
		g.GET("/client-profile", h.GetClientProfile)
		g.GET("/dian-config", h.GetDianConfig)
		g.PUT("/dian-config", h.SaveDianConfig)
	}
}
